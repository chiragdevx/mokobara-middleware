package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleProductRequest(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	fmt.Printf("🔥 Received event: %+v\n", request)
	fmt.Printf("🔥 Request body: %s\n", request.Body)

	shopifyTopic := request.Headers["X-Shopify-Topic"]
	fmt.Printf("🔥 X-Shopify-Topic header: %s\n", shopifyTopic)

	var data map[string]interface{}
	err := json.Unmarshal([]byte(request.Body), &data)
	if err != nil {
		fmt.Printf("❌ Error unmarshalling body: %v\n", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf(`{"error": "Invalid JSON: %v"}`, err),
		}, nil
	}

	switch shopifyTopic {
	case "products/create":
		fmt.Println("🔥 Handling 'products/create' event")
		// createProductHandler(request)

	case "products/update":
		fmt.Println("🔥 Handling 'products/update' event")
		updateProductHandler(data)
	default:
		fmt.Println("🔥 Unknown Shopify Topic:", shopifyTopic)
	}

	fmt.Println("✅")
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{"message": "test"}`,
	}, nil

}

func updateProductHandler(data map[string]interface{}) {
	productId := data["id"].(float64)
	productIdStr := fmt.Sprintf("%.0f", productId)

	product, err := getProductWithMetafields(productIdStr)
	if err != nil {
		fmt.Printf("❌ Error fetching product: %v\n", err)
		return
	}

	payload, err := getProductPayload(product)
	if err != nil {
		fmt.Printf("❌ Error generating payload: %v\n", err)
		return
	}

	fmt.Printf("🔥 Main Payload: %+v\n", payload)

	var wg sync.WaitGroup
	errChan := make(chan error, len(payload)) // Capture errors from goroutines

	for _, product := range payload {
		wg.Add(1)
		go func(product map[string]interface{}) {
			defer wg.Done()
			if err := createProduct(product); err != nil {
				errChan <- err
			}
		}(product)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		fmt.Printf("❌ API call failed : %v\n", err)
	}
}

func main() {
	lambda.Start(HandleProductRequest)
}
