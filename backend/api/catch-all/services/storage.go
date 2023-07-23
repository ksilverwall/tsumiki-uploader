package services

import (
	"catch-all/repositories"
	"errors"
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
)

var (
	ErrUnexpected          = errors.New("unexpected server error")
	ErrThumbnailNotCreated = errors.New("thumbnail has not been created")
)

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

func GenerateID() (string, error) {
	u7, err := uuid.NewV7()
	if err != nil {
		ErrorLog(fmt.Errorf("failed to create transaction id: %w", err).Error())
		return "", ErrUnexpected
	}

	id := u7.String()

	return id, nil
}

type Storage struct {
	StorageRepository repositories.Storage
}

func (s *Storage) GetFileUploadUrl(id string) (string, string, error) {
	filePath := fmt.Sprintf("%v.zip", id)

	url, err := s.StorageRepository.GetSignedUrl(repositories.SignedUrlModePUT, filePath)
	if err != nil {
		ErrorLog(fmt.Errorf("failed to get signed url: %w", err).Error())
		return "", "", ErrUnexpected
	}

	return url, filePath, nil
}

func (s *Storage) GetFileDownloadUrl(id string) (string, error) {
	url, err := s.StorageRepository.GetSignedUrl(repositories.SignedUrlModeGET, fmt.Sprintf("%v.zip", id))
	if err != nil {
		ErrorLog(err.Error())
		return "", ErrUnexpected
	}

	return url, nil
}

func (s *Storage) GetFileThumbnailUrls(key string) ([]string, error) {
	data, err := s.StorageRepository.Get(fmt.Sprintf("thumbnails/%v/_keys", key))
	if err != nil {
		ErrorLog(err.Error())
		return []string{}, ErrThumbnailNotCreated
	}

	if len(data) == 0 {
		ErrorLog("thumbnail key data is empty")
		return []string{}, ErrUnexpected
	}
	keys := SplitLines(string(data))
	if len(keys) == 0 {
		ErrorLog("thumbnail key file is empty")
		return []string{}, ErrUnexpected
	}

	urls := make([]string, len(keys))
	for i, key := range keys {
		url, err := s.StorageRepository.GetSignedUrl(repositories.SignedUrlModeGET, key)
		if err != nil {
			ErrorLog(fmt.Errorf("failed to create signed url: %w", err).Error())
			return []string{}, ErrUnexpected
		}
		urls[i] = url
	}

	if len(urls) == 0 {
		ErrorLog("items is empty")
		return []string{}, ErrUnexpected
	}

	return urls, nil
}
