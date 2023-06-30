package main

import (
	"errors"

	"catch-all/apigateway"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
	MainRouter.Register(Server{})
	lambda.Start(handler)
}
