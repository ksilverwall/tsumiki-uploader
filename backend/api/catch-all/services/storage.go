package services

import (
	"catch-all/models"
	"catch-all/repositories"
	"fmt"
	"strings"
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

func GetArchivePath(s models.FileID) string {
	return fmt.Sprintf("files/%v/archive.zip", s)
}

func GetThumbnailsPrefix(s models.FileID) string {
	return fmt.Sprintf("files/%v/thumbnails/", s)
}

func GetThumbnailsKeyPath(s models.FileID) string {
	return fmt.Sprintf("files/%v/thumbnails/_keys", s)
}

type Storage struct {
	StorageRepository repositories.Storage
}

func (s *Storage) GetFileUploadUrl(id models.FileID) (string, string, error) {
	filePath := GetArchivePath(id)

	url, err := s.StorageRepository.GetSignedUrl(repositories.SignedUrlModePUT, filePath)
	if err != nil {
		ErrorLog(fmt.Errorf("failed to get signed url: %w", err).Error())
		return "", "", ErrUnexpected
	}

	InfoLog(fmt.Sprintf("get upload url for id '%s'", id))

	return url, filePath, nil
}

func (s *Storage) GetFileDownloadUrl(id models.FileID) (string, error) {
	url, err := s.StorageRepository.GetSignedUrl(repositories.SignedUrlModeGET, GetArchivePath(id))
	if err != nil {
		ErrorLog(err.Error())
		return "", ErrUnexpected
	}

	InfoLog(fmt.Sprintf("get download url for id '%s'", id))

	return url, nil
}

func (s *Storage) GetFileThumbnailUrls(id models.FileID) ([]string, error) {
	data, err := s.StorageRepository.Get(GetThumbnailsKeyPath(id))
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
