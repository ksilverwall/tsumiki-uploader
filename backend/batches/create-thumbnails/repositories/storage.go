package repositories

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

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

func (r StorageRepository) Upload(key string, filePath string) error {
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

func (r StorageRepository) UploadThumbnails(dirpath string, key string, filePaths []string) error {
	prefix := fmt.Sprintf("thumbnails/%s", key)

	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("file not found")
	}

	s3paths := []string{}
	for _, file := range filePaths {
		s3path := filepath.Join(prefix, filepath.Base(file))

		err = r.Upload(s3path, file)
		if err != nil {
			return err
		}

		s3paths = append(s3paths, s3path)
	}

	keysFile, err := ioutil.TempFile("", "*")
	if err != nil {
		return fmt.Errorf("Failed to create key file: %w", err)
	}

	for _, s3p := range s3paths {
		_, err = keysFile.WriteString(s3p)
		if err != nil {
			return err
		}
	}

	keysFile.Close()

	err = r.Upload(filepath.Join(prefix, "_keys"), keysFile.Name())
	if err != nil {
		return err
	}

	return nil
}
