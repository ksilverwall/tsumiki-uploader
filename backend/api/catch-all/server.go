package main

import (
	"catch-all/gen/openapi"
	"catch-all/models"
	"catch-all/repositories"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type Server struct {
	StorageRepository      repositories.Storage
	TransactionRepository  repositories.Transaction
	StateMachineRepository repositories.StateMachine
}

func (s Server) GetFileUrl(ctx *gin.Context, key string) {
	url, err := s.StorageRepository.GetSignedUrl(repositories.SignedUrlModeGET, fmt.Sprintf("%v.zip", key))
	if err != nil {
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, openapi.Error{Code: http.StatusInternalServerError, Message: "failed to get signed url"})
			return
		}
	}

	ctx.JSON(http.StatusOK, openapi.DownloadInfo{
		Name: "dummy_name.zip",
		Url:  url,
	})
}

func (s Server) GetFileThumbnailUrls(ctx *gin.Context, key string) {
	data, err := s.StorageRepository.Get(fmt.Sprintf("thumbnails/%v/_keys", key))
	if err != nil {
		ctx.JSON(http.StatusNotFound, openapi.Error{Code: http.StatusNotFound, Message: "thumbnail is not created"})
		return
	}

	if len(data) == 0 {
		ctx.JSON(http.StatusInternalServerError, openapi.Error{Code: http.StatusInternalServerError, Message: "thumbnail key data is empty"})
		return
	}
	keys := SplitLines(string(data))
	if len(keys) == 0 {
		ctx.JSON(http.StatusInternalServerError, openapi.Error{Code: http.StatusInternalServerError, Message: "thumbnail key file is empty"})
		return
	}

	items := []openapi.FileThumbnail{}
	for _, key := range keys {
		url, err := s.StorageRepository.GetSignedUrl(repositories.SignedUrlModeGET, key)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, openapi.Error{Code: http.StatusInternalServerError, Message: fmt.Errorf("failed to create signed url: %w", err).Error()})
			return
		}

		items = append(items, openapi.FileThumbnail{Url: url})
	}

	if len(items) == 0 {
		ctx.JSON(http.StatusInternalServerError, openapi.Error{Code: http.StatusInternalServerError, Message: "items is empty"})
		return
	}

	ctx.JSON(http.StatusOK, openapi.FileThumbnails{
		Items: items,
	})
}

func (s Server) CreateTransaction(ctx *gin.Context) {
	u7, err := uuid.NewV7()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, openapi.Error{Code: http.StatusInternalServerError, Message: "failed to create transaction id"})
		return
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

	url, err := s.StorageRepository.GetSignedUrl(repositories.SignedUrlModePUT, filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, openapi.Error{Code: http.StatusInternalServerError, Message: "failed to get signed url"})
		return
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

func SplitLines(data string) []string {
	ret := []string{}
	for _, str := range strings.Split(data, "\n") {
		buf := strings.TrimSpace(str)
		if len(buf) > 0 {
			ret = append(ret, buf)
		}
	}

	return ret
}
