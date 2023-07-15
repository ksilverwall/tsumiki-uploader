package main

import (
	"catch-all/gen/openapi"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

type Server struct {
	AWSSession *session.Session
	BucketName string
}

func (s Server) CreateTransaction(ctx echo.Context) error {
	u7, err := uuid.NewV7()
	svc := s3.New(s.AWSSession)
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(fmt.Sprintf("%v.zip", u7)),
	})
	url, err := req.Presign(15 * time.Minute)
	if err != nil {
		panic(err)
	}

	return ctx.JSON(http.StatusOK, openapi.Transaction{
		Id:  u7.String(),
		Url: url,
	})
}

func (s Server) GetFileUrl(ctx echo.Context, key string) error {
	svc := s3.New(s.AWSSession)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(fmt.Sprintf("%v.zip", key)),
	})
	url, err := req.Presign(15 * time.Minute)
	if err != nil {
		panic(err)
	}

	return ctx.JSON(http.StatusOK, openapi.DownloadInfo{
		Name: "dummy_name.zip",
		Url:  url,
	})
}

func (s Server) GetFileThumbnailUrls(ctx echo.Context, key string) error {
	svc := s3.New(s.AWSSession)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(fmt.Sprintf("thumnails/%v/thum-0", key)),
	})
	url, err := req.Presign(15 * time.Minute)
	if err != nil {
		panic(err)
	}

	return ctx.JSON(http.StatusOK, openapi.FileThumbnails{
		Items: []openapi.FileThumbnail{
			{
				Url: url,
			},
		},
	})
}
