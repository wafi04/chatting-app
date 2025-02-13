package comments

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	authrepository "github.com/wafi04/chatting-app/services/auth/pkg/repository"
	"github.com/wafi04/chatting-app/services/shared/types"
	"github.com/wafi04/chatting-app/services/shared/utils"
)

type CommentRepository struct {
	db       *sqlx.DB
	authrepo *authrepository.AuthRepository
}

func NewCommentRepository(db *sqlx.DB, authrepo *authrepository.AuthRepository) *CommentRepository {
	return &CommentRepository{
		db:       db,
		authrepo: authrepo,
	}
}

func (r *CommentRepository) CreateComment(ctx context.Context, req *types.CreateComment) (*types.Comment, error) {
	// Validate PostID
	if req.PostID == "" {
		return nil, fmt.Errorf("PostID cannot be empty")
	}

	// Check if the post exists
	var postExists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)", req.PostID).Scan(&postExists)
	if err != nil {
		return nil, fmt.Errorf("failed to check post existence: %w", err)
	}
	if !postExists {
		return nil, fmt.Errorf("post with ID %s does not exist", req.PostID)
	}

	// Handle ParentID and calculate depth
	var parentID interface{} = nil // Default to NULL for top-level comments
	var depth int                  // Depth of the new comment

	if req.ParentID != nil && *req.ParentID != "" {
		var parentExists bool
		var parentDepth int
		err := r.db.QueryRowContext(ctx, `
            SELECT EXISTS(SELECT 1 FROM comments WHERE id = $1), depth 
            FROM comments WHERE id = $1
        `, *req.ParentID).Scan(&parentExists, &parentDepth)
		if err != nil {
			return nil, fmt.Errorf("failed to check parent comment existence: %w", err)
		}
		if !parentExists {
			return nil, fmt.Errorf("parent comment with ID %s does not exist", *req.ParentID)
		}
		parentID = *req.ParentID
		depth = parentDepth + 1 // Increment depth based on the parent's depth
	} else {
		// Top-level comment
		depth = 0
	}

	// Insert the new comment into the database
	query := `
        INSERT INTO comments (id, user_id, post_id, content, depth, created_at, parent_comment_id) 
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, user_id, post_id, content, depth, created_at, parent_comment_id
    `
	comment := &types.Comment{}
	err = r.db.QueryRowContext(
		ctx,
		query,
		utils.GenerateRandomId("coment"),
		req.UserID,
		req.PostID,
		req.Content,
		depth, // Use the calculated depth
		time.Now(),
		parentID,
	).Scan(
		&comment.ID,
		&comment.UserID,
		&comment.PostID,
		&comment.Content,
		&comment.Depth,
		&comment.CreatedAT,
		&comment.ParentID,
	)
	if err != nil {
		return nil, err
	}

	// Fetch user info
	user, err := r.authrepo.GetUser(ctx, &types.GetUserRequest{
		UserId: comment.UserID,
	})
	if err != nil {
		return nil, err
	}
	comment.UserInfo = *user

	return comment, nil
}
func (r *CommentRepository) GetParentDepth(ctx context.Context, parentID string) (int32, error) {
	var depth int32
	err := r.db.QueryRowContext(ctx, "SELECT depth FROM comments WHERE id = $1", parentID).Scan(&depth)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("parent category not found")
		}
		return 0, fmt.Errorf("failed to get parent category: %v", err)
	}
	return depth, nil
}
func (r *CommentRepository) GetCommentTree(ctx context.Context, postID string) (map[string]*types.Comment, []*types.Comment, error) {
	query := `
        WITH RECURSIVE comment_tree AS (
            SELECT 
                c.id, 
                c.user_id, 
                c.post_id, 
                c.content, 
                c.parent_comment_id, 
                c.depth, 
                c.created_at,
                ARRAY[c.id]::VARCHAR[] AS path,
                0 AS level
            FROM comments c
            WHERE c.parent_comment_id IS NULL AND c.post_id = $1
            UNION ALL
            SELECT 
                c.id, 
                c.user_id, 
                c.post_id, 
                c.content, 
                c.parent_comment_id, 
                c.depth, 
                c.created_at,
                ct.path || c.id::VARCHAR,
                ct.level + 1
            FROM comments c
            INNER JOIN comment_tree ct ON ct.id = c.parent_comment_id
        )
        SELECT 
            id, 
            user_id, 
            post_id, 
            content, 
            depth, 
            created_at, 
            parent_comment_id,
            path
        FROM comment_tree
        ORDER BY path
    `
	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query comments: %v", err)
	}
	defer rows.Close()

	commentMap := make(map[string]*types.Comment)
	var rootComments []*types.Comment

	for rows.Next() {
		var comm types.Comment
		var createdAt sql.NullTime
		var parentID sql.NullString
		var path pq.StringArray

		err := rows.Scan(
			&comm.ID,
			&comm.UserID,
			&comm.PostID,
			&comm.Content,
			&comm.Depth,
			&createdAt,
			&parentID,
			&path,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan comment: %v", err)
		}

		// Set nullable fields
		if parentID.Valid {
			comm.ParentID = &parentID.String
		}
		if createdAt.Valid {
			comm.CreatedAT = createdAt.Time
		}

		user, _ := r.authrepo.GetUser(ctx, &types.GetUserRequest{
			UserId: comm.UserID,
		})

		comm.Path = path
		comm.UserInfo = *user

		// Add to comment map
		commentMap[comm.ID] = &comm

		if !parentID.Valid {
			rootComments = append(rootComments, &comm)
		} else {
			parent := commentMap[parentID.String]
			if parent != nil {
				if parent.Replies == nil {
					parent.Replies = []*types.Comment{}
				}
				parent.Replies = append(parent.Replies, &comm)
			}
		}
	}

	if err = rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error iterating comments: %v", err)
	}

	return commentMap, rootComments, nil
}
func (s *CommentRepository) DeleteComment(ctx context.Context, req *types.DeleteComment) (*types.DeleteCommentReponse, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Validate comment existence
	var exists bool
	if err = tx.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM comments WHERE id = $1)",
		req.CommentID,
	).Scan(&exists); err != nil {
		return nil, fmt.Errorf("failed to check comment existence: %w", err)
	}
	if !exists {
		return nil, &CommentNotFoundError{CommentID: req.CommentID}
	}

	var deletedCount int64
	if req.DeleteChildren {
		// Delete comment and all descendants using recursive CTE
		deletedCount, err = s.deleteCommentTree(ctx, tx, req.CommentID)
	} else {
		// Delete single comment and update children
		deletedCount, err = s.deleteSingleComment(ctx, tx, req.CommentID)
	}
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &types.DeleteCommentReponse{
		Success: true,
		Count:   deletedCount,
	}, nil
}

type CommentNotFoundError struct {
	CommentID string
}

func (e *CommentNotFoundError) Error() string {
	return fmt.Sprintf("comment with ID %s not found", e.CommentID)
}

// deleteCommentTree removes a comment and all its descendants
func (s *CommentRepository) deleteCommentTree(ctx context.Context, tx *sql.Tx, commentID string) (int64, error) {
	const query = `
        WITH RECURSIVE comment_tree AS (
            SELECT id 
            FROM comments 
            WHERE id = $1
            UNION ALL
            SELECT c.id
            FROM comments c
            INNER JOIN comment_tree ct ON c.parent_comment_id = ct.id
        )
        DELETE FROM comments
        WHERE id IN (SELECT id FROM comment_tree)
        RETURNING id`

	result, err := tx.ExecContext(ctx, query, commentID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete comment tree: %w", err)
	}

	return result.RowsAffected()
}

// deleteSingleComment removes a single comment and updates its children
func (s *CommentRepository) deleteSingleComment(ctx context.Context, tx *sql.Tx, commentID string) (int64, error) {
	// Check for children
	var hasChildren bool
	if err := tx.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM comments WHERE parent_comment_id = $1)",
		commentID,
	).Scan(&hasChildren); err != nil {
		return 0, fmt.Errorf("failed to check for children: %w", err)
	}

	if hasChildren {
		// Update children to have no parent
		if _, err := tx.ExecContext(ctx,
			"UPDATE comments SET parent_comment_id = NULL WHERE parent_comment_id = $1",
			commentID,
		); err != nil {
			return 0, fmt.Errorf("failed to update children: %w", err)
		}
	}

	// Delete the comment
	result, err := tx.ExecContext(ctx,
		"DELETE FROM comments WHERE id = $1",
		commentID,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to delete comment: %w", err)
	}

	return result.RowsAffected()
}

type RespCount struct {
	Count int64 `json:"count"`
}

func (s *CommentRepository) GetCountComment(ctx context.Context, postID string) (*RespCount, error) {
	var count int64 // Ubah ke int64 untuk menghindari masalah dengan tipe data COUNT
	query := `
    SELECT COUNT(*)
    FROM comments
    WHERE post_id = $1
    `
	err := s.db.QueryRowContext(ctx, query, postID).Scan(&count)

	if err != nil {
		return nil, fmt.Errorf("failed to get count: %w", err)
	}

	return &RespCount{
		Count: count, // Konversi ke int32 jika RespCount.Count bertipe int32
	}, nil
}
