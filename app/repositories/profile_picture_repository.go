package repositories

import (
	"context"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

type ProfilePictureRepository interface {
	BucketExists(bucketName string) bool
	CreateBucket(bucketName string) error
	CreateProfilePicture(file *multipart.FileHeader, bucketName, name string) error
}

func NewProfilePictureRepository(minioClient *minio.Client) ProfilePictureRepository {
	return &profilePictureRepository{ctx: context.Background(), minioClient: minioClient}
}

type profilePictureRepository struct {
	ctx         context.Context
	minioClient *minio.Client
}

func (r *profilePictureRepository) BucketExists(bucketName string) bool {
	exists, _ := r.minioClient.BucketExists(r.ctx, bucketName)
	return exists
}

func (r *profilePictureRepository) CreateBucket(bucketName string) error {
	return r.minioClient.MakeBucket(r.ctx, bucketName, minio.MakeBucketOptions{})
}

func (r *profilePictureRepository) CreateProfilePicture(file *multipart.FileHeader, bucketName, fileName string) error {
	srcFile, err := file.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()

	_, err = r.minioClient.PutObject(r.ctx, bucketName, fileName, srcFile, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return err
	}

	return nil
}
