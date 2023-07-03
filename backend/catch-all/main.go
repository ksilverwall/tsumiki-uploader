package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"catch-all/apigateway"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	MainRouter  = apigateway.NewRouter()
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

	ctx := apigateway.NewContext()
	route := MainRouter.GetRoute(request.Path, request.HTTPMethod)
	if route == nil {
		return events.APIGatewayProxyResponse{}, errors.New("Route not found")
	}

	if err := route.Handler(ctx); err != nil {
		// ERROR: Invalid format for parameter key: parameter 'key' is empty, can't bind its value
		return events.APIGatewayProxyResponse{}, err
	}

	if ctx.Status.Response == nil {
		return events.APIGatewayProxyResponse{}, errors.New("response not found")
	}

	return events.APIGatewayProxyResponse{
		Headers:    CORSHeaders,
		StatusCode: ctx.Status.Response.StatusCode,
		Body:       ctx.Status.Response.Body,
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
	MainRouter.Register(Server{
		AWSSession: sess,
		BucketName: os.Getenv("STORAGE_BUCKET_NAME"),
	})
	lambda.Start(handler)
}
