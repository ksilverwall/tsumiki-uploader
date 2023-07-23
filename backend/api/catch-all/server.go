package main

import (
	"catch-all/gen/openapi"
	"catch-all/models"
	"catch-all/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServerError struct {
	Status int
	Body   openapi.Error
}

func NewServerError(err error) ServerError {
	if errors.Is(err, services.ErrThumbnailNotCreated) {
		return ServerError{
			Status: http.StatusNotFound,
			Body:   openapi.Error{Code: http.StatusNotFound, Message: err.Error()},
		}
	}
	if errors.Is(err, services.ErrUnexpected) {
		return ServerError{
			Status: http.StatusInternalServerError,
			Body:   openapi.Error{Code: http.StatusInternalServerError, Message: err.Error()},
		}
	}

	return ServerError{
		Status: http.StatusInternalServerError,
		Body:   openapi.Error{Code: http.StatusInternalServerError, Message: err.Error()},
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
		ctx.JSON(http.StatusInternalServerError, openapi.Error{Code: http.StatusInternalServerError, Message: err.Error()})
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
		se := NewServerError(err)
		ctx.JSON(se.Status, se.Body)
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
		se := NewServerError(err)
		ctx.JSON(se.Status, se.Body)
		return
	}

	url, filePath, err := s.StorageService.GetFileUploadUrl(id)
	if err != nil {
		se := NewServerError(err)
		ctx.JSON(se.Status, se.Body)
		return
	}

	err = s.TransactionService.Create(id, filePath)
	if err != nil {
		se := NewServerError(err)
		ctx.JSON(se.Status, se.Body)
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
		se := NewServerError(err)
		ctx.JSON(se.Status, se.Body)
		return
	}

	err = s.AsyncService.CreateThumbnails(models.ThumbnailRequest{
		TransactionID: transactionId,
		FilePath:      t.FilePath,
	})
	if err != nil {
		se := NewServerError(err)
		ctx.JSON(se.Status, se.Body)
	}

	ctx.JSON(http.StatusOK, openapi.Transaction{
		Id:     transactionId,
		Status: openapi.UPLOADED,
	})
}
