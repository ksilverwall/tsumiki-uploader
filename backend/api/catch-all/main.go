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

type GlobalParams struct {
	GinEngine *gin.Engine
}

var (
	Global = GlobalParams{
		GinEngine: gin.Default(),
	}
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

	Global.GinEngine.ServeHTTP(resp, req)

	return events.APIGatewayProxyResponse{
		Headers:    CORSHeaders,
		StatusCode: resp.Code,
		Body:       resp.Body.String(),
	}, nil
}

func Init(g *GlobalParams) error {
	region := os.Getenv("STORAGE_REGION")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}

	params, err := GetParameters(sess)
	if err != nil {
		return fmt.Errorf("failed to get parameters: %v", err)
	}

	server := CreateServer(params)
	if err != nil {
		return fmt.Errorf("failed to init server: %v", err)
	}

	openapi.RegisterHandlers(g.GinEngine, server)

	return nil
}

func main() {
	err := Init(&Global)
	if err != nil {
		log.Println(fmt.Sprintf("ERROR: %s", err.Error()))
		return
	}
	lambda.Start(handler)
}
