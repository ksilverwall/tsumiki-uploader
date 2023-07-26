package main

import (
	"catch-all/gen/openapi"
	"catch-all/models"
	"catch-all/services"
	"errors"
	"net/http"
	"time"

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
	if e, ok := err.(*services.Error); ok {
		switch e.Code {
		case services.ErrLabelThumbnailNotCreated:
			return ErrorResponse{
				Status: http.StatusNotFound,
				Body:   openapi.ClientError{Code: openapi.ClientErrorCodeThumbnailNotFound, Message: "thumbnail has not been created"},
			}
		default:
			return ErrorResponse{
				Status: http.StatusInternalServerError,
				Body:   openapi.ClientError{Code: openapi.ClientErrorCodeUnknown, Message: "unexpected server error"},
			}
		}
	}
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

func (s Server) GetFileUrl(ctx *gin.Context, id string) {
	url, err := s.StorageService.GetFileDownloadUrl(models.FileID(id))
	if err != nil {
		NewServerError(err).Send(ctx)
		return
	}

	ctx.JSON(http.StatusOK, openapi.DownloadInfo{
		Name: "download.zip",
		Url:  url,
	})
}

func (s Server) GetFileThumbnailUrls(ctx *gin.Context, id string) {
	urls, err := s.StorageService.GetFileThumbnailUrls(models.FileID(id))
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
	v, exists := ctx.Get("RequestTimeSec")
	if !exists {
		services.ErrorLog("RequestTimeSec is not exists")
		NewServerError(services.ErrUnexpected).Send(ctx)
		return
	}

	currentTime, ok := v.(time.Time)
	if !ok {
		services.ErrorLog("RequestTimeSec is not time obj")
		NewServerError(services.ErrUnexpected).Send(ctx)
		return
	}

	id, err := services.GenerateID()
	if err != nil {
		NewServerError(err).Send(ctx)
		return
	}

	fid := models.FileID(id)

	url, filePath, err := s.StorageService.GetFileUploadUrl(fid)
	if err != nil {
		NewServerError(err).Send(ctx)
		return
	}

	tid, err := s.TransactionService.Create(fid, filePath, currentTime)
	if err != nil {
		NewServerError(err).Send(ctx)
		return
	}

	ctx.JSON(http.StatusOK, openapi.Transaction{
		Id:     string(tid),
		Status: openapi.READY,
		Url:    url,
	})
}

func (s Server) UpdateTransaction(ctx *gin.Context, transactionId string) {
	tid := models.TransactionID(transactionId)
	t, err := s.TransactionService.Get(tid)
	if err != nil {
		NewServerError(err).Send(ctx)
		return
	}

	err = s.AsyncService.CreateThumbnails(models.ThumbnailRequest{
		TransactionID:         tid,
		ArchiveFilePath:       t.FilePath,
		ThumbnailFilesKeyPath: services.GetThumbnailsKeyPath(t.FileID),
		ThumbnailFilesPrefix:  services.GetThumbnailsPrefix(t.FileID),
	})
	if err != nil {
		NewServerError(err).Send(ctx)
		return
	}

	ctx.JSON(http.StatusOK, openapi.Transaction{
		Id:     transactionId,
		Status: openapi.UPLOADED,
		FileId: string(t.FileID),
	})
}
