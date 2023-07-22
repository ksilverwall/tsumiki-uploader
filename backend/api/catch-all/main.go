package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"catch-all/gen/openapi"
	"catch-all/models"
	"catch-all/repositories"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/gin-gonic/gin"
)

var (
	InitError   error
	GinEngine   = gin.Default()
	CORSHeaders = map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type, Authorization",
		"Access-Control-Allow-Methods": "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod == http.MethodOptions {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Headers:    CORSHeaders,
		}, nil
	}

	req, _ := http.NewRequest(request.HTTPMethod, request.Path, strings.NewReader(request.Body))
	resp := httptest.NewRecorder()

	GinEngine.ServeHTTP(resp, req)

	return events.APIGatewayProxyResponse{
		Headers:    CORSHeaders,
		StatusCode: resp.Code,
		Body:       resp.Body.String(),
	}, nil
}

func uninitHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Headers:    CORSHeaders,
		StatusCode: http.StatusInternalServerError,
		Body:       fmt.Sprintf("Failed to init server: %v", InitError),
	}, nil
}

func NewServer(region string) (openapi.ServerInterface, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Println(fmt.Sprintf("failed to create session: %v", err))
		return Server{}, nil
	}

	var pp models.PlatformParameters
	var pb models.BatchParameters

	parameterRepository := repositories.ParameterRepository{Client: ssm.New(sess)}

	err = parameterRepository.Get("/app/tsumiki-uploader/backend/platform", &pp)
	if err != nil {
		return Server{}, fmt.Errorf("failed to load platform parameters: %w", err)
	}

	err = parameterRepository.Get("/app/tsumiki-uploader/backend/batches", &pb)
	if err != nil {
		return Server{}, fmt.Errorf("failed to load batch parameters: %w", err)
	}

	server := Server{
		AWSSession: sess,
		BucketName: pp.DataStorage,
		TransactionRepository: repositories.Transaction{
			Dynamodb:  dynamodb.New(sess),
			TableName: pp.TransactionTable.Name,
		},
		StateMachineRepository: repositories.StateMachine{
			Client:          sfn.New(sess),
			StateMachineArn: pb.ThumbnailsCreatingStateMachineArn,
		},
	}

	return server, nil
}

func main() {
	server, err := NewServer(os.Getenv("STORAGE_REGION"))
	if err != nil {
		log.Println(fmt.Sprintf("failed to init server: %v", err))
		InitError = err
		return
	}

	openapi.RegisterHandlers(GinEngine, server)

	lambda.Start(handler)
}
