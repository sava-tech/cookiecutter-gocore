package cloudinary

import (
	"context"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	Cld *cloudinary.Cloudinary
}

func NewCloudinaryService() (*CloudinaryService, error) {
	cld, err := cloudinary.NewFromParams(
		"djhxnetly",
		"158867669949168",
		"Z004YVSE9e6ZQJblqw8UyYVqukQ",
	)
	if err != nil {
		return nil, err
	}

	return &CloudinaryService{Cld: cld}, nil
}

func (cs *CloudinaryService) UploadImage(ctx context.Context, file interface{}, folder string) (string, string, error) {
	result, err := cs.Cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder: folder,
	})
	if err != nil {
		return "", "", err
	}

	return result.SecureURL, result.PublicID, nil
}
