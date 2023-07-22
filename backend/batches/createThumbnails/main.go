package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

type InEvent struct {
	Payload string `json:"payload"`
}

func handler(request InEvent) (InEvent, error) {
	return request, nil
}

func main() {
	lambda.Start(handler)
}
