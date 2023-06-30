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
	u2, err := uuid.NewV7()
	svc := s3.New(s.AWSSession)
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(fmt.Sprintf("%v.zip", u2)),
	})
	url, err := req.Presign(15 * time.Minute)
	if err != nil {
		panic(err)
	}

	return ctx.JSON(http.StatusOK, openapi.Transaction{
		Id:  0,
		Url: url,
	})
}
