package s3

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	s3uploader "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	s3Client *s3.Client
	bucket   string
)

func InitS3() error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(getEnv("AWS_REGION", "us-east-1")))
	if err != nil {
		panic("failed to load AWS config: " + err.Error())
	}
	s3Client = s3.NewFromConfig(cfg)
	bucket = getEnv("AWS_BUCKET_NAME", "")
	return nil
}

func UploadProfilePicture(file multipart.File, fileHeader *multipart.FileHeader, userID int) (string, error) {
	defer file.Close()

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(file)
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("profile_pictures/user_%d_%d%s", userID, time.Now().Unix(), filepath.Ext(fileHeader.Filename))

	uploader := s3uploader.NewUploader(s3Client)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(filename),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", err
	}

	return result.Location, nil
}

func getEnv(key, fallback string) string {
	if val, ok := syscall.Getenv(key); ok {
		return val
	}
	return fallback
}
