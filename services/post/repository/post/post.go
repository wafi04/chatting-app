package postrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	authrepository "github.com/wafi04/chatting-app/services/auth/pkg/repository"
	"github.com/wafi04/chatting-app/services/comments"
	"github.com/wafi04/chatting-app/services/shared/pkg/logger"
	"github.com/wafi04/chatting-app/services/shared/types"
	"github.com/wafi04/chatting-app/services/shared/utils"
)

type PostRepository struct {
	DB          *sqlx.DB
	logger      logger.Logger
	commentRepo *comments.CommentRepository
	authrepo    *authrepository.AuthRepository
}

func NewPostRepository(db *sqlx.DB, commentRepo *comments.CommentRepository, authrepo *authrepository.AuthRepository) *PostRepository {
	return &PostRepository{
		DB:          db,
		commentRepo: commentRepo,
		authrepo:    authrepo,
	}
}
func (r *PostRepository) CreatePost(ctx context.Context, req *types.Post) (*types.Post, error) {
	tx, err := r.DB.BeginTxx(ctx, nil)
	if err != nil {
		r.logger.Log(logger.ErrorLevel, "Failed to begin transaction: %v", err)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				r.logger.Log(logger.ErrorLevel, "Failed to commit transaction: %v", err)
			}
		}
	}()

	ID := utils.GenerateRandomId("POST")

	tags := pq.Array(req.Tags)
	mentions := pq.Array(req.Mentions)

	var post types.Post
	var created_at time.Time
	var dbTags, dbMentions []string

	query := `
    INSERT INTO posts (
        id,
        user_id,
        caption,
        location,
        tags,
        mentions,
        created_at
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING 
        id,
        user_id,
        caption,
        location,
        tags,
        mentions,
        created_at
    `

	err = tx.QueryRowContext(ctx, query,
		ID,
		req.UserId,
		req.Caption,
		req.Location,
		tags,
		mentions,
		time.Now(),
	).Scan(
		&post.Id,
		&post.UserId,
		&post.Caption,
		&post.Location,
		pq.Array(&dbTags),
		pq.Array(&dbMentions),
		&created_at,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	post.Tags = dbTags
	post.Mentions = dbMentions

	for _, mediaUpload := range req.Media {
		mediaID := utils.GenerateRandomId("MEDIA")
		r.logger.Log(logger.InfoLevel, "url : %s", mediaUpload.FileUrl)
		_, err := r.CreateMedia(ctx, tx, &types.Media{
			Id:       mediaID,
			FileUrl:  mediaUpload.FileUrl,
			PublicId: mediaUpload.PublicId,
			FileType: mediaUpload.FileType,
			FileName: mediaUpload.FileName,
			PostId:   post.Id,
		})
		if err != nil {
			r.logger.Log(logger.ErrorLevel, "Failed to upload media: %v", err)
			return nil, fmt.Errorf("failed to upload media: %v", err)
		}
	}

	post.CreatedAt = created_at.Unix()
	post.Media = req.Media
	r.logger.Log(logger.InfoLevel, "res: id=%s, user_id=%s, caption=%s", post.Id, post.UserId, post.Caption)
	return &post, nil
}
func (r *PostRepository) GetUserPosts(ctx context.Context, req *types.GetUserPostsRequest) (*types.GetUserPostsResponse, error) {
	query := `
    SELECT 
        id,
        user_id,
        caption,
        location,
        tags,
        mentions,
        created_at,
        updated_at
    FROM 
        posts
    WHERE 
        user_id = $1
    ORDER BY 
        created_at DESC
    LIMIT $2 OFFSET $3
    `

	limit := req.PerPage
	offset := (req.Page - 1) * req.PerPage

	posts, err := r.QueryPosts(ctx, query, req.UserId, limit, offset)
	if err != nil {
		return nil, err
	}

	return &types.GetUserPostsResponse{
		Posts: posts,
	}, nil
}

func (r *PostRepository) GetAllPosts(ctx context.Context, req *types.GetAllPostsRequest) (*types.GetAllPostsResponse, error) {
	query := `
    SELECT 
        id,
        user_id,
        caption,
        location,
        tags,
        mentions,
        created_at,
        updated_at
    FROM 
        posts
    ORDER BY 
        created_at DESC
    LIMIT $1 OFFSET $2
    `

	limit := req.PerPage
	offset := (req.Page - 1) * req.PerPage

	posts, err := r.QueryPosts(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return &types.GetAllPostsResponse{
		Posts: posts,
	}, nil
}

func (r *PostRepository) DeletePosts(ctx context.Context, req *types.DeletePostRequest) (*types.DeletePostResponse, error) {
	query := `
        DELETE FROM posts 
        WHERE id = $1
        RETURNING id`

	result, err := r.DB.ExecContext(ctx, query, req.PostId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete post: %v", err)
	}

	// Cek apakah ada row yang terpengaruh
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("post with id %s not found", req.PostId)
	}

	return &types.DeletePostResponse{
		Success: true,
	}, nil
}
