package postservice

import (
	"context"
	"fmt"

	cloudrepo "github.com/wafi04/chatting-app/services/post/pkg/repository/cloud"
	postrepo "github.com/wafi04/chatting-app/services/post/pkg/repository/post"
	"github.com/wafi04/chatting-app/services/shared/pkg/logger"
	"github.com/wafi04/chatting-app/services/shared/types"
)

type PostService struct {
	cloudrepo *cloudrepo.Cloudinary
	postrepo  *postrepo.PostRepository
	logger    logger.Logger
}

func NewPostService(
	cloudrepo *cloudrepo.Cloudinary,
	postrepo *postrepo.PostRepository,
) *PostService {
	return &PostService{
		cloudrepo: cloudrepo,
		postrepo:  postrepo,
	}
}

func (s *PostService) CreatePost(ctx context.Context, req *types.CreatePostRequest) (*types.PostResponse, error) {
	s.logger.Log(logger.InfoLevel, "Incoming create post request from user: %s", req.UserId)

	var uploadedMedia []*types.Media
	for _, mediaUpload := range req.Media {
		media, err := s.cloudrepo.UploadFile(ctx, &types.MediaUpload{
			FileData: mediaUpload.FileData,
			FileName: mediaUpload.FileName,
			FileType: mediaUpload.FileType,
		})
		if err != nil {
			s.logger.Log(logger.ErrorLevel, "Failed to upload media: %v", err)
			s.rollbackMediaUploads(ctx, uploadedMedia)
			return nil, fmt.Errorf("failed to upload media: %v", err)
		}
		uploadedMedia = append(uploadedMedia, media)
	}

	post, err := s.postrepo.CreatePost(ctx, &types.Post{
		UserId:   req.UserId,
		Caption:  req.Caption,
		Media:    uploadedMedia,
		Mentions: req.Mentions,
		Location: req.Location,
		Tags:     req.Tags,
	})
	if err != nil {
		s.logger.Log(logger.ErrorLevel, "Failed to create post: %v", err)
		s.rollbackMediaUploads(ctx, uploadedMedia)
		return nil, fmt.Errorf("failed to create post: %v", err)
	}

	return &types.PostResponse{
		Post: post,
	}, nil
}

func (s *PostService) rollbackMediaUploads(ctx context.Context, media []*types.Media) {
	for _, m := range media {
		if err := s.cloudrepo.DeleteFile(ctx, m.PublicId); err != nil {
			s.logger.Log(logger.ErrorLevel, "Failed to delete media during rollback: %v", err)
		}
	}
}

func (s *PostService) GetUserPosts(ctx context.Context, req *types.GetUserPostsRequest) (*types.GetUserPostsResponse, error) {
	return s.postrepo.GetUserPosts(ctx, req)
}

func (s *PostService) GetAllPosts(ctx context.Context, req *types.GetAllPostsRequest) (*types.GetAllPostsResponse, error) {
	return s.postrepo.GetAllPosts(ctx, req)
}
