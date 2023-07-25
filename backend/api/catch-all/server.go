package main

import (
	"catch-all/gen/openapi"
	"catch-all/models"
	"catch-all/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Status int
	Body   any
}

func (r ErrorResponse) Send(ctx *gin.Context) {
	ctx.JSON(r.Status, r.Body)
}

func NewServerError(err error) ErrorResponse {
	if errors.Is(err, services.ErrThumbnailNotCreated) {
		return ErrorResponse{
			Status: http.StatusNotFound,
			Body:   openapi.ClientError{Code: openapi.ClientErrorCodeThumbnailNotFound, Message: err.Error()},
		}
	}
	if errors.Is(err, services.ErrUnexpected) {
		return ErrorResponse{
			Status: http.StatusInternalServerError,
			Body:   openapi.ClientError{Code: openapi.ClientErrorCodeUnknown, Message: err.Error()},
		}
	}

	return ErrorResponse{
		Status: http.StatusInternalServerError,
		Body:   openapi.ServerError{Code: openapi.ServerErrorCodeUnknown, Message: err.Error()},
	}
}

type Server struct {
	StorageService     services.Storage
	TransactionService services.Transaction
	AsyncService       services.Async
}

func (s Server) GetFileUrl(ctx *gin.Context, key string) {
	url, err := s.StorageService.GetFileDownloadUrl(key)
	if err != nil {
		NewServerError(err).Send(ctx)
		return
	}

	ctx.JSON(http.StatusOK, openapi.DownloadInfo{
		Name: "download.zip",
		Url:  url,
	})
}

func (s Server) GetFileThumbnailUrls(ctx *gin.Context, key string) {
	urls, err := s.StorageService.GetFileThumbnailUrls(key)
	if err != nil {
		NewServerError(err).Send(ctx)
		return
	}

	items := make([]openapi.FileThumbnail, len(urls))
	for i, u := range urls {
		items[i] = openapi.FileThumbnail{Url: u}
	}

	ctx.JSON(http.StatusOK, openapi.FileThumbnails{
		Items: items,
	})
}

func (s Server) CreateTransaction(ctx *gin.Context) {
	id, err := services.GenerateID()
	if err != nil {
		NewServerError(err).Send(ctx)
		return
	}

	url, filePath, err := s.StorageService.GetFileUploadUrl(id)
	if err != nil {
		NewServerError(err).Send(ctx)
		return
	}

	err = s.TransactionService.Create(id, filePath)
	if err != nil {
		NewServerError(err).Send(ctx)
		return
	}

	ctx.JSON(http.StatusOK, openapi.Transaction{
		Id:     id,
		Status: openapi.READY,
		Url:    url,
	})
}

func (s Server) UpdateTransaction(ctx *gin.Context, transactionId string) {
	t, err := s.TransactionService.Get(transactionId)
	if err != nil {
		NewServerError(err).Send(ctx)
		return
	}

	err = s.AsyncService.CreateThumbnails(models.ThumbnailRequest{
		TransactionID: transactionId,
		FilePath:      t.FilePath,
	})
	if err != nil {
		NewServerError(err).Send(ctx)
		return
	}

	ctx.JSON(http.StatusOK, openapi.Transaction{
		Id:     transactionId,
		Status: openapi.UPLOADED,
	})
}
