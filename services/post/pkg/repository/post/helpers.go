package postrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/wafi04/chatting-app/services/shared/types"
)

func (r *PostRepository) QueryPosts(ctx context.Context, query string, args ...interface{}) ([]*types.Post, error) {
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}
	defer rows.Close()

	var posts []*types.Post
	for rows.Next() {
		post := &types.Post{}
		var dbTags, dbMentions []string
		var created_at, updated_at time.Time

		err := rows.Scan(
			&post.Id,
			&post.UserId,
			&post.Caption,
			&post.Location,
			pq.Array(&dbTags),
			pq.Array(&dbMentions),
			&created_at,
			&updated_at,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}

		post.Tags = dbTags
		post.Mentions = dbMentions
		post.CreatedAt = created_at.Unix()
		post.UpdatedAt = updated_at.Unix()

		mediaList, err := r.GetMediaByPost(ctx, nil, post.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to get media for post %s: %w", post.Id, err)
		}

		post.Media = mediaList

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating rows: %w", err)
	}

	return posts, nil
}
