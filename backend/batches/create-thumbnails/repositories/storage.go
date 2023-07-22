package repositories

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type StorageRepository struct {
	Downloader *s3manager.Downloader
	Uploader   *s3manager.Uploader
	BucketName string
}

func (r StorageRepository) GetZip(key string) ([]byte, error) {
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

func (r StorageRepository) UploadThumbnails(dirpath string, key string) error {
	prefix := fmt.Sprintf("thumbnails/%s", key)

	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("file not found")
	}

	keyFile, err := ioutil.TempFile("", "*")
	if err != nil {
		return fmt.Errorf("Failed to create key file: %w", err)
	}
	defer keyFile.Close()

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(dirpath, file.Name())
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		s3path := filepath.Join(prefix, strings.TrimPrefix(filePath, dirpath))

		_, err = r.Uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(r.BucketName),
			Key:    aws.String(s3path),
			Body:   file,
		})
		if err != nil {
			return err
		}

		_, err = keyFile.WriteString(s3path)
		if err != nil {
			return err
		}
	}

	_, err = r.Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(r.BucketName),
		Key:    aws.String(filepath.Join(prefix, "_keys")),
		Body:   keyFile,
	})
	if err != nil {
		return err
	}

	return nil
}
