package repositories

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
)

type SignedUrlMode int

const (
	SignedUrlModeGET SignedUrlMode = iota
	SignedUrlModePUT
)

type Storage struct {
	Client     *s3.S3
	BucketName string
}

func (r Storage) Get(key string) ([]byte, error) {
	obj, err := r.Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(r.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return []byte{}, err
	}

	bodyBytes, err := ioutil.ReadAll(obj.Body)
	if err != nil {
		return []byte{}, err
	}

	return bodyBytes, nil
}

func (r Storage) GetSignedUrl(mode SignedUrlMode, key string) (string, error) {
	var req *request.Request

	if mode == SignedUrlModeGET {
		req, _ = r.Client.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(r.BucketName),
			Key:    aws.String(key),
		})
	}
	if mode == SignedUrlModePUT {
		req, _ = r.Client.PutObjectRequest(&s3.PutObjectInput{
			Bucket: aws.String(r.BucketName),
			Key:    aws.String(key),
		})
	} else {
		return "", fmt.Errorf("unsupported signed url mode")
	}

	url, err := req.Presign(15 * time.Minute)
	if err != nil {
		return "", fmt.Errorf("failed to presign: %w", err)
	}

	return url, nil
}
