package main

import (
	"encoding/json"
	"fmt"
	"log"

	"create-thumbnails/models"

	"github.com/aws/aws-lambda-go/lambda"
)

type StateMachineRequest struct {
	Input models.ThumbnailRequest
}

func handler(request StateMachineRequest) (models.ThumbnailRequest, error) {
	b, _ := json.Marshal(request)
	log.Println(fmt.Sprintf("INFO: START with request '%s'", string(b)))
	err := MainProcess(request.Input)
	if err != nil {
		log.Println(fmt.Sprintf("ERROR: %s", err.Error()))
	}

	return models.ThumbnailRequest{}, err
}

func main() {
	lambda.Start(handler)
}
