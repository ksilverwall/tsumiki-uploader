package main

import (
	"catch-all/gen/openapi"
	"catch-all/models"
	"catch-all/repositories"
	"catch-all/services"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type ServerIntegrationParameters struct {
	Session  *session.Session
	Platform models.PlatformParameters
	Batch    models.BatchParameters
}

func GetParameters(sess *session.Session) (ServerIntegrationParameters, error) {
	var params ServerIntegrationParameters
	var err error

	params.Session = sess

	parameterRepository := repositories.ParameterRepository{Client: ssm.New(sess)}

	err = parameterRepository.Get("/app/tsumiki-uploader/backend/platform", &params.Platform)
	if err != nil {
		return params, fmt.Errorf("failed to load platform parameters: %w", err)
	}

	err = parameterRepository.Get("/app/tsumiki-uploader/backend/batches", &params.Batch)
	if err != nil {
		return params, fmt.Errorf("failed to load batch parameters: %w", err)
	}

	return params, err
}

func CreateServer(params ServerIntegrationParameters) openapi.ServerInterface {
	server := Server{
		StorageService: services.Storage{
			StorageRepository: repositories.Storage{
				Client:     s3.New(params.Session),
				BucketName: params.Platform.DataStorage,
			},
		},
		TransactionService: services.Transaction{
			TransactionRepository: repositories.Transaction{
				Dynamodb:  dynamodb.New(params.Session),
				TableName: params.Platform.TransactionTable.Name,
			},
		},
		AsyncService: services.Async{
			StateMachineRepository: repositories.StateMachine{
				Client:          sfn.New(params.Session),
				StateMachineArn: params.Batch.ThumbnailsCreatingStateMachineArn,
			},
		},
	}

	return server
}
