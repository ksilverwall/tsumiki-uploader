package main

import (
	"catch-all/gen/openapi"
	"catch-all/models"
	"catch-all/repositories"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type Server struct {
	AWSSession             *session.Session
	BucketName             string
	TransactionRepository  repositories.Transaction
	StateMachineRepository repositories.StateMachine
}

func (s Server) GetFileUrl(ctx *gin.Context, key string) {
	svc := s3.New(s.AWSSession)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(fmt.Sprintf("%v.zip", key)),
	})
	url, err := req.Presign(15 * time.Minute)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, openapi.DownloadInfo{
		Name: "dummy_name.zip",
		Url:  url,
	})
}

func (s Server) GetFileThumbnailUrls(ctx *gin.Context, key string) {
	svc := s3.New(s.AWSSession)

	urls, err := LoadThumbnailPaths(svc, s.BucketName, key)
	if err != nil {
		ctx.JSON(http.StatusNotFound, openapi.Error{Code: http.StatusNotFound, Message: "thumbnail is not created"})
		return
	}

	items := make([]openapi.FileThumbnail, len(urls))
	for i, url := range urls {
		req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(s.BucketName),
			Key:    aws.String(url),
		})
		url, err := req.Presign(15 * time.Minute)
		if err != nil {
			panic(err)
		}

		items[i] = openapi.FileThumbnail{Url: url}
	}

	ctx.JSON(http.StatusOK, openapi.FileThumbnails{
		Items: items,
	})
}

func (s Server) CreateTransaction(ctx *gin.Context) {
	var err error
	u7, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}

	id := u7.String()
	filePath := fmt.Sprintf("%v.zip", id)

	err = s.TransactionRepository.Put(id, models.Transaction{
		ID:       id,
		FilePath: filePath,
	})
	if err != nil {
		errRes := NewErrorResponse(fmt.Errorf("failed to push transaction: %w", err))
		ctx.JSON(errRes.Status, errRes.Body)
		return
	}

	svc := s3.New(s.AWSSession)
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(filePath),
	})
	url, err := req.Presign(15 * time.Minute)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, openapi.Transaction{
		Id:     id,
		Status: openapi.READY,
		Url:    url,
	})
}

func (s Server) UpdateTransaction(ctx *gin.Context, transactionId string) {
	t, err := s.TransactionRepository.Get(transactionId)
	if err != nil {
		e := fmt.Errorf("transaction id not found: %w", err)
		ctx.JSON(http.StatusNotFound, openapi.Error{Code: http.StatusNotFound, Message: e.Error()})
		return
	}

	err = s.StateMachineRepository.Execute(models.ThumbnailRequest{
		TransactionID: transactionId,
		FilePath:      t.FilePath,
	})
	if err != nil {
		e := fmt.Errorf("transaction id not found: %w", err)
		ctx.JSON(http.StatusInternalServerError, openapi.Error{
			Code:    http.StatusInternalServerError,
			Message: e.Error(),
		})
	}

	ctx.JSON(http.StatusOK, openapi.Transaction{
		Id:     transactionId,
		Status: openapi.UPLOADED,
	})
}

func LoadThumbnailPaths(svc *s3.S3, bucketName, key string) ([]string, error) {
	obj, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fmt.Sprintf("thumnails/%v/_keys", key)),
	})
	if err != nil {
		return []string{}, err
	}

	bodyBytes, err := ioutil.ReadAll(obj.Body)
	if err != nil {
		return []string{}, err
	}

	ret := []string{}
	for _, str := range strings.Split(string(bodyBytes), "\n") {
		buf := strings.TrimSpace(str)
		if len(buf) > 0 {
			ret = append(ret, buf)
		}
	}

	return ret, nil
}
