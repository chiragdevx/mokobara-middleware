package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// HandleOrderRequest handles API Gateway requests
func HandleOrderRequest(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	fmt.Printf("🔥 Received event: %+v\n", request)
	fmt.Printf("🔥 Request body: %s\n", request.Body)

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{"message": "test"}`,
	}, nil
}

func main() {
	lambda.Start(HandleOrderRequest)
}
