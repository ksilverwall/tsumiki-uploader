package main

import (
	"create-thumbnails/repositories"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

func CreateKeyFile(s3paths []string) (string, error) {
	keysFile, err := ioutil.TempFile("", "*")
	if err != nil {
		return "", fmt.Errorf("Failed to create key file: %w", err)
	}

	for _, s3p := range s3paths {
		_, err = keysFile.WriteString(s3p)
		if err != nil {
			return "", err
		}
	}

	keysFile.Close()

	return keysFile.Name(), nil
}

type Service struct {
	StoraStorageRepository repositories.Storage
}

func (s Service) UploadThumbnails(keyFilePath string, prefix string, filePaths []string) error {
	s3paths := []string{}
	for _, file := range filePaths {
		s3path := filepath.Join(prefix, filepath.Base(file))

		err := s.StoraStorageRepository.Upload(s3path, file)
		if err != nil {
			return err
		}

		s3paths = append(s3paths, s3path)
	}

	kfn, err := CreateKeyFile(s3paths)
	if err != nil {
		return nil
	}

	err = s.StoraStorageRepository.Upload(keyFilePath, kfn)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) CreateThumbnails(archiveFilePath string, tempDir string) ([]string, error) {
	reader, err := s.StoraStorageRepository.GetZip(archiveFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get zip file in '%s': %w", archiveFilePath, err)
	}

	result, err := ProcessZipFiles(reader, tempDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create thumbnails: %w", err)
	}

	for _, err := range result.Errs {
		log.Println(fmt.Errorf("WARNING: %w", err))
	}
	if len(result.FilePaths) == 0 {
		return nil, fmt.Errorf("no thumbnail files created")
	}

	return result.FilePaths, nil
}
