package main

import (
	"archive/zip"
	"bytes"
	"create-thumbnails/models"
	"create-thumbnails/repositories"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/nfnt/resize"
)

type Result struct {
	FilePaths []string
	Errs      []error
}

func ProcessZipFile(f *zip.File, outPath string) (string, error) {
	if f.FileInfo().IsDir() {
		return "", fmt.Errorf("directory detected")
	}

	rc, err := f.Open()
	buffer := make([]byte, 512)
	rc.Read(buffer)
	contentType := http.DetectContentType(buffer)

	rc, err = f.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer rc.Close()

	img, _, err := image.Decode(rc)
	if err != nil {
		return "", fmt.Errorf("failed to decode image file (type %s): %w", contentType, err)
	}

	err = CreateThumbnail(img, outPath)
	if err != nil {
		return "", fmt.Errorf("failed to create thumbnail file: %w", err)
	}

	return outPath, nil
}

func ProcessZipFiles(content []byte, outDir string) (Result, error) {
	errs := []error{}
	paths := []string{}
	reader, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return Result{}, err
	}

	if len(reader.File) == 0 {
		return Result{}, fmt.Errorf("File not found: %w", err)
	}

	for i, f := range reader.File {
		_, fileName := path.Split(f.Name)
		outPath := path.Join(outDir, strings.Replace(fileName, ".", "_thumbnail.", 1))

		fname, err := ProcessZipFile(f, outPath)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to process file[%d]: %w", i, err))
			continue
		}
		paths = append(paths, fname)
	}

	return Result{FilePaths: paths, Errs: errs}, nil
}

func CreateThumbnail(img image.Image, outPath string) error {
	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	thumbnail := resize.Thumbnail(100, 100, img, resize.Lanczos3)

	err = jpeg.Encode(outFile, thumbnail, nil)
	if err != nil {
		return fmt.Errorf("failed to encode: %w", err)
	}

	return nil
}

func MainProcess(request models.ThumbnailRequest) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION"))},
	)

	if err != nil {
		return fmt.Errorf("failed to init session: %w", err)
	}

	var pp models.PlatformParameters

	parameterRepository := repositories.ParameterRepository{Client: ssm.New(sess)}

	err = parameterRepository.Get("/app/tsumiki-uploader/platform/storages/backend", &pp)
	if err != nil {
		return fmt.Errorf("failed to load platform parameters: %w", err)
	}

	s := Service{
		StoraStorageRepository: repositories.Storage{
			Downloader: s3manager.NewDownloader(sess),
			Uploader:   s3manager.NewUploader(sess),
			BucketName: pp.DataStorage,
		},
	}

	tempDir, err := os.MkdirTemp("", "temp_*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	paths, err := s.CreateThumbnails(request.ArchiveFilePath, tempDir)
	if err != nil {
		return fmt.Errorf("failed to create thumbnails: %w", err)
	}

	err = s.UploadThumbnails(request.ThumbnailFilesKeyPath, request.ThumbnailFilesPrefix, paths)
	if err != nil {
		return fmt.Errorf("failed to upload thumbnails: %w", err)
	}

	return nil
}
