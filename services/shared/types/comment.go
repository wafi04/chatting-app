package types

import (
	"time"

	"github.com/lib/pq"
)

type CreateComment struct {
	PostID    string    `json:"postId"`
	UserID    string    `json:"userId"`
	Content   string    `json:"content"`
	ParentID  *string   `json:"parentId,omitempty"`
	CreatedAT time.Time `json:"createdAt"`
	Depth     int64     `json:"depth"`
}

type Comment struct {
	ID        string         `db:"id" json:"Id"`
	PostID    string         `db:"post_id" json:"postId"`
	UserID    string         `db:"user_id" json:"userId"`
	Content   string         `db:"content" json:"content"`
	Depth     int64          `db:"depth" json:"depth"`
	Path      pq.StringArray `json:"-"`
	UserInfo  UserInfo       `json:"user_info"`
	CreatedAT time.Time      `db:"created_at" json:"createdAt"`
	Replies   []*Comment     `json:"replies"`
	ParentID  *string        `db:"parent_id" json:"parentId"`
}

type DeleteComment struct {
	CommentID      string `json:"commentID"`
	DeleteChildren bool   `json:"delete"`
}

type DeleteCommentReponse struct {
	Success bool  `json:"success"`
	Count   int64 `json:"count"`
}

type ListCommentsRequest struct {
	Page            int32   `json:"page"`
	Limit           int32   `json:"limit"`
	PostID          string  `json:"post_id"`
	ParentID        *string `json:"parent_id,omitempty"`
	IncludeChildren bool    `json:"include_children"`
}
type ListCommentsResponse struct {
	Comments []*Comment `json:"comments"`
	Total    int32      `json:"total"`
}
