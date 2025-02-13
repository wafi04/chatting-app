package comments_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/wafi04/chatting-app/services/shared/pkg/constants"
	"github.com/wafi04/chatting-app/services/shared/types"
)

func TestCreateComment(t *testing.T) {
	fixedTime := time.Now()

	tests := []struct {
		name    string
		input   *types.CreateComment
		mockFn  func(mock sqlmock.Sqlmock)
		want    *types.Comment
		wantErr bool
	}{
		{
			name: "successful comment creation",
			input: &types.CreateComment{
				PostID:  "POST1234",
				UserID:  "uuexkbkabaka",
				Content: "djdkddkd",
				Depth:   0,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "post_id", "content", "depth", "created_at", "parent_comment_id",
				}).AddRow(
					constants.CommentID,
					"uuexkbkabaka",
					"POST1234",
					"djdkddkd",
					0,
					fixedTime,
					"isssusysys",
				)

				mock.ExpectQuery(`INSERT INTO comments \(id,user_id,post_id,content,depth,created_at,parent_comment_id\)\s+VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7\)\s+RETURNING\s+id,user_id,post_id,content,depth,created_at,parent_comment_id`).
					WithArgs(
						constants.CommentID,
						"uuexkbkabaka",
						"POST1234",
						"djdkddkd",
						0,
						sqlmock.AnyArg(),
						"isssusysys",
					).
					WillReturnRows(rows)
			},
			want: &types.Comment{
				ID:        constants.CommentID,
				UserID:    "uuexkbkabaka",
				PostID:    "POST1234",
				Content:   "djdkddkd",
				Depth:     0,
				CreatedAT: fixedTime,
			},
			wantErr: false,
		},
		{
			name: "successful comment creation without parent",
			input: &types.CreateComment{
				PostID:  "POST1234",
				UserID:  "uuexkbkabaka",
				Content: "djdkddkd",
				Depth:   0,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "post_id", "content", "depth", "created_at", "parent_comment_id",
				}).AddRow(
					constants.CommentID,
					"uuexkbkabaka",
					"POST1234",
					"djdkddkd",
					0,
					fixedTime,
					"",
				)

				mock.ExpectQuery(`INSERT INTO comments \(id,user_id,post_id,content,depth,created_at,parent_comment_id\)\s+VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7\)\s+RETURNING\s+id,user_id,post_id,content,depth,created_at,parent_comment_id`).
					WithArgs(
						constants.CommentID,
						"uuexkbkabaka",
						"POST1234",
						"djdkddkd",
						0,
						sqlmock.AnyArg(),
						"",
					).
					WillReturnRows(rows)
			},
			want: &types.Comment{
				ID:        constants.CommentID,
				UserID:    "uuexkbkabaka",
				PostID:    "POST1234",
				Content:   "djdkddkd",
				Depth:     0,
				CreatedAT: fixedTime,
			},
			wantErr: false,
		},
		{
			name: "database error",
			input: &types.CreateComment{
				PostID:  "POST1234",
				UserID:  "uuexkbkabaka",
				Content: "djdkddkd",
				Depth:   0,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO comments \(id,user_id,post_id,content,depth,created_at,parent_comment_id\)\s+VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7\)\s+RETURNING\s+id,user_id,post_id,content,depth,created_at,parent_comment_id`).
					WithArgs(
						constants.CommentID,
						"uuexkbkabaka",
						"POST1234",
						"djdkddkd",
						0,
						sqlmock.AnyArg(),
					).
					WillReturnError(sqlmock.ErrCancelled)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock database connection
			mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
			if err != nil {
				t.Fatalf("Failed to create mock: %v", err)
			}
			defer mockDB.Close()

			// Wrap the mock database with sqlx
			db := sqlx.NewDb(mockDB, "sqlmock")
			defer db.Close()

			// Initialize repository with mock db
			// repo := comments.NewCommentRepository(db)

			// Setup mock expectations
			tt.mockFn(mock)

			// Execute the function
			// got, err := repo.CreateComment(context.Background(), tt.input)

			// Check error
			// if tt.wantErr {
			// 	assert.Error(t, err)
			// 	assert.Nil(t, got)
			// 	return
			// }

			// assert.NoError(t, err)
			// assert.NotNil(t, got)

			// // Check result
			// assert.Equal(t, tt.want.ID, got.ID)
			// assert.Equal(t, tt.want.UserID, got.UserID)
			// assert.Equal(t, tt.want.PostID, got.PostID)
			// assert.Equal(t, tt.want.Content, got.Content)
			// assert.Equal(t, tt.want.Depth, got.Depth)
			// assert.Equal(t, tt.want.ParentID, got.ParentID)

			// // Ensure all expectations were met
			// if err := mock.ExpectationsWereMet(); err != nil {
			// 	t.Errorf("there were unfulfilled expectations: %s", err)
			// }
		})
	}
}
