package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"catch-all/gen/openapi"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
)

var (
	GinEngine   = gin.Default()
	CORSHeaders = map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type, Authorization",
		"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
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

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("STORAGE_REGION")),
	})
	if err != nil {
		log.Println(fmt.Sprintf("failed to create session: %v", err))
		return
	}

	openapi.RegisterHandlers(GinEngine, Server{
		AWSSession: sess,
		BucketName: os.Getenv("STORAGE_BUCKET_NAME"),
	})

	lambda.Start(handler)
}
