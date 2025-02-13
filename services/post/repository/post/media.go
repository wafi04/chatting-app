package postrepo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wafi04/chatting-app/services/shared/types"
)

func (r *PostRepository) CreateMedia(ctx context.Context, tx *sqlx.Tx, req *types.Media) (*types.Media, error) {
	var media types.Media
	var created_at time.Time
	query := `
	INSERT INTO media
	(
		id,
		post_id,
		file_url,
		public_id,
		file_type,
		file_name,
		created_at
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7)
	RETURNING
		id,
		post_id,
		file_url,
		public_id,
		file_type,
		file_name,
		created_at
	`

	err := tx.QueryRowContext(ctx, query, req.Id,
		req.PostId,
		req.FileUrl,
		req.PublicId,
		req.FileType,
		req.FileName,
		time.Now(),
	).Scan(
		&media.Id,
		&media.PostId,
		&media.FileUrl,
		&media.PublicId,
		&media.FileType,
		&media.FileName,
		&created_at,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create post : %w", err)
	}

	return &media, nil
}

func (r *PostRepository) GetMediaByPost(ctx context.Context, tx *sqlx.Tx, postId string) ([]*types.Media, error) {
	query := `
    SELECT 
        id,
        file_url,
        public_id,
        file_type,
        file_name,
        created_at
    FROM 
        media
    WHERE 
        post_id = $1
    `

	var rows *sql.Rows
	var err error
	if tx != nil {
		rows, err = tx.QueryContext(ctx, query, postId)
	} else {
		rows, err = r.DB.QueryContext(ctx, query, postId)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get media for post %s: %w", postId, err)
	}
	defer rows.Close()

	var mediaList []*types.Media
	for rows.Next() {
		media := &types.Media{}
		var mediaCreatedAt time.Time

		err := rows.Scan(
			&media.Id,
			&media.FileUrl,
			&media.PublicId,
			&media.FileType,
			&media.FileName,
			&mediaCreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan media for post %s: %w", postId, err)
		}

		media.CreatedAt = mediaCreatedAt.Unix()
		media.PostId = postId
		mediaList = append(mediaList, media)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating media rows: %w", err)
	}

	return mediaList, nil
}
