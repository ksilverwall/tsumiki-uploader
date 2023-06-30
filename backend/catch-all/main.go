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
	route := MainRouter.GetRoute(request.Path, request.HTTPMethod)
	if route == nil {
		return events.APIGatewayProxyResponse{}, errors.New("Route not found")
	}

	res := events.APIGatewayProxyResponse{}

	if err := route.Handler(apigateway.NewContext(&res)); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return res, nil
}

func main() {
	MainRouter.Register(Server{})
	lambda.Start(handler)
}
