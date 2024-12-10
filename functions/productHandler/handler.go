package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

// HandleProductRequest handles API Gateway requests
func HandleProductRequest(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	fmt.Println("🔥 Hello, world")
	response := events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{"message": "Hello, world"}`,
	}
	return response, nil
}
