package comments

import (
	"context"

	"github.com/wafi04/chatting-app/services/shared/types"
)

type CommentService struct {
	repo *CommentRepository
}

func NewCommntService(repo *CommentRepository) *CommentService {
	return &CommentService{
		repo: repo,
	}
}

func (s *CommentService) CreateComment(ctx context.Context, req *types.CreateComment) (*types.Comment, error) {
	return s.repo.CreateComment(ctx, req)
}
func (s *CommentService) GetComments(ctx context.Context, req *types.ListCommentsRequest) (*types.ListCommentsResponse, error) {
	categoryMap, rootCategories, err := s.repo.GetCommentTree(ctx, req.PostID)
	if err != nil {
		return nil, err
	}

	return &types.ListCommentsResponse{
		Comments: rootCategories,
		Total:    int32(len(categoryMap)),
	}, nil
}
func (s *CommentService) DeleteComment(ctx context.Context, req *types.DeleteComment) (*types.DeleteCommentReponse, error) {
	return s.repo.DeleteComment(ctx, req)
}
