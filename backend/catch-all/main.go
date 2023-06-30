package main

import (
	"errors"
	"fmt"
	"log"

	"catch-all/apigateway"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	MainRouter = apigateway.NewRouter()
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx := apigateway.NewContext()
	route := MainRouter.GetRoute(request.Path, request.HTTPMethod)
	if route == nil {
		return events.APIGatewayProxyResponse{}, errors.New("Route not found")
	}

	if err := route.Handler(ctx); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if ctx.Status.Response == nil {
		return events.APIGatewayProxyResponse{}, errors.New("response not found")
	}

	return *ctx.Status.Response, nil
}

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		log.Println(fmt.Sprintf("failed to create session: %v", err))
		return
	}
	MainRouter.Register(Server{
		AWSSession: sess,
		BucketName: "dummy-bucket-name",
	})
	lambda.Start(handler)
}
