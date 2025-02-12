package cloudrepo

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/wafi04/chatting-app/services/shared/types"
	"github.com/wafi04/chatting-app/services/shared/utils"
)

type Cloudinary struct {
	cloudinary *cloudinary.Cloudinary
}

func NewCloudinaryService(cld *cloudinary.Cloudinary) *Cloudinary {
	return &Cloudinary{cloudinary: cld}
}

func (s *Cloudinary) UploadFile(
	ctx context.Context,
	req *types.MediaUpload,
) (*types.Media, error) {
	FileID := utils.GenerateRandomId("files")
	tempFile, err := ioutil.TempFile("", "upload-*")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write(req.FileData); err != nil {
		return nil, err
	}

	tempFile.Close()

	uploadResult, err := s.cloudinary.Upload.Upload(ctx, tempFile.Name(), uploader.UploadParams{
		Folder: "ChattingApp",
		PublicID: utils.GenerateCustomID(utils.IDOptions{
			Prefix:       "Media",
			CustomFormat: "{prefix}_{rand:6}_{timestamp}",
		}),
	})
	if err != nil {
		return nil, err
	}

	return &types.Media{
		Id:       FileID,
		FileUrl:  uploadResult.URL,
		PublicId: uploadResult.PublicID,
		FileType: req.FileType,
		FileName: req.FileName,
	}, nil
}

func (s *Cloudinary) DeleteFile(ctx context.Context, publicID string) error {
	_, err := s.cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from cloudinary: %v", err)
	}
	return nil
}
