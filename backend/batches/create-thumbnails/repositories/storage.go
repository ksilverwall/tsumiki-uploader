package repositories

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Storage struct {
	Downloader *s3manager.Downloader
	Uploader   *s3manager.Uploader
	BucketName string
}

func (r Storage) GetZip(key string) ([]byte, error) {
	buff := aws.WriteAtBuffer{}

	_, err := r.Downloader.Download(&buff,
		&s3.GetObjectInput{
			Bucket: aws.String(r.BucketName),
			Key:    aws.String(key),
		})
	if err != nil {
		return []byte{}, fmt.Errorf("failed to download file '%s': %w", key, err)
	}

	return buff.Bytes(), nil
}

func (r Storage) Upload(key string, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = r.Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(r.BucketName),
		Key:    aws.String(key),
		Body:   file,
	})
	return err
}
