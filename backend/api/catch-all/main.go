package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"catch-all/gen/openapi"
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
	if InitError != nil {
		return events.APIGatewayProxyResponse{
			Headers:    CORSHeaders,
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Failed to init server: %v", InitError),
		}, nil
	}

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

func initServer() (openapi.ServerInterface, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("STORAGE_REGION")),
	})
	if err != nil {
		log.Println(fmt.Sprintf("failed to create session: %v", err))
		return Server{}, nil
	}

	svc := ssm.New(sess)
	paramsInput := &ssm.GetParametersByPathInput{
		Path:           aws.String("/app/tsumiki-uploader/backend"),
		WithDecryption: aws.Bool(true),
		Recursive:      aws.Bool(true),
	}

	params, err := svc.GetParametersByPath(paramsInput)
	if err != nil {
		return Server{}, fmt.Errorf("failed to load parameter: %w", err)
	}

	paramsMap := make(map[string]string)
	for _, param := range params.Parameters {
		paramsMap[*param.Name] = *param.Value
	}

	stateMachineArn := paramsMap["/app/tsumiki-uploader/backend/batches/thumbnails-creating-state-machine"]
	if len(stateMachineArn) == 0 {
		return Server{}, errors.New("stateMachineArn is not set")
	}

	TableName := paramsMap["/app/tsumiki-uploader/backend/transaction-table/name"]
	if len(TableName) == 0 {
		return Server{}, errors.New("table name is not set")
	}

	server := Server{
		AWSSession: sess,
		BucketName: os.Getenv("STORAGE_BUCKET_NAME"),
		TransactionRepository: repositories.Transaction{
			Dynamodb:  dynamodb.New(sess),
			TableName: TableName,
		},
		StateMachineRepository: repositories.StateMachine{
			Client:          sfn.New(sess),
			StateMachineArn: stateMachineArn,
		},
	}

	return server, nil
}

func main() {
	server, err := initServer()
	if err != nil {
		log.Println(fmt.Sprintf("failed to init server: %v", err))
		InitError = err
		return
	}
	openapi.RegisterHandlers(GinEngine, server)

	lambda.Start(handler)
}
