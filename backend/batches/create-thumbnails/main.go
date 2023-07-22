package main

import (
	"fmt"
	"log"
	"os"

	"create-thumbnails/models"
	"create-thumbnails/repositories"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type ThumbnailRequest struct {
	TransactionID string
	FilePath      string
}

type StateMachineRequest struct {
	Input ThumbnailRequest
}

func mainProcess(request ThumbnailRequest) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION"))},
	)

	if err != nil {
		return fmt.Errorf("failed to init session: %w", err)
	}

	var pp models.PlatformParameters

	parameterRepository := repositories.ParameterRepository{Client: ssm.New(sess)}

	err = parameterRepository.Get("/app/tsumiki-uploader/backend/platform", &pp)
	if err != nil {
		return fmt.Errorf("failed to load platform parameters: %w", err)
	}

	r := repositories.StorageRepository{
		Downloader: s3manager.NewDownloader(sess),
		Uploader:   s3manager.NewUploader(sess),
		BucketName: pp.DataStorage,
	}

	reader, err := r.GetZip(request.FilePath)
	if err != nil {
		return fmt.Errorf("failed to get zip file: %w", err)
	}

	tempDir, err := os.MkdirTemp("", "temp_*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	result, err := ProcessZipFiles(reader, tempDir)
	if err != nil {
		return fmt.Errorf("failed to create thumbnails: %w", err)
	}

	for _, err := range result.Errs {
		log.Println(fmt.Errorf("WARNING: %w", err))
	}
	if len(result.FilePaths) == 0 {
		return fmt.Errorf("no thumbnail files created")
	}

	err = r.UploadThumbnails(request.TransactionID, result.FilePaths)
	if err != nil {
		return fmt.Errorf("failed to upload thumbnails: %w", err)
	}

	return nil
}

func handler(request StateMachineRequest) (ThumbnailRequest, error) {
	err := mainProcess(request.Input)
	if err != nil {
		log.Println(fmt.Errorf("ERROR: %w", err))
	}

	return ThumbnailRequest{}, err
}

func main() {
	lambda.Start(handler)
}
