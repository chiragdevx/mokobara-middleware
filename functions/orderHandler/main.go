package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// HandleProductRequest handles API Gateway requests
func HandleProductRequest(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// Log the incoming event and request body
	fmt.Printf("🔥 Received event: %+v\n", request)
	fmt.Printf("🔥 Request body: %s\n", request.Body)

	// Parse the request body into a map (assuming JSON format)
	var requestBody map[string]interface{}
	if err := json.Unmarshal([]byte(request.Body), &requestBody); err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       `{"message": "Invalid request body"}`,
		}, nil
	}

	// Get product payload from the request body
	payload := getProductPayload(requestBody)
	fmt.Println("payload:", payload)


	// Construct a response with the payload
	response := events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: fmt.Sprintf(`{"product": %+v}`, payload),
	}

	return response, nil
}

// getProductPayload processes the request body to construct a product payload
func getProductPayload(body map[string]interface{}) string {
	// Prepare the custom attributes array
	

	// Construct the payload
	productPayload := map[string]interface{}{
		"product": map[string]interface{}{
			"sku":       "your-product-sku", 
			"name":      body["title"],     
			"attribute_set_id": 4,
			"price":         99.99, 
			"status":        1, 
			"visibility":    4, 
			"type_id":       "simple",
			"weight":        1.0, 
			"extension_attributes": map[string]interface{}{
				"stock_item": map[string]interface{}{
					"qty":      10,   
					"is_in_stock": true,
				},
			},
				
		},
	}

	// Convert the payload to JSON format
	payloadJSON, err := json.Marshal(productPayload)
	if err != nil {
		fmt.Printf("Error marshalling product payload: %v\n", err)
		return `{"error": "Internal server error"}`
	}

	// Return the JSON payload as a string
	return string(payloadJSON)
}

func main() {
	lambda.Start(HandleProductRequest)
}
