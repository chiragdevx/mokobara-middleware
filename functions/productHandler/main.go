package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"


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
	products := getProductsPayload(requestBody)

	


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
func getProductsPayload(body map[string]interface{}) string {
	products := []map[string]interface{}{}

	// Safely assert variants to []interface{}
	variants, ok := body["variants"].([]interface{})
	if !ok {
		fmt.Println("⚠️ 'variants' is not a slice of interface{}")
		return "[]"
	}

	for _, item := range variants {
		// Safely assert each variant as map[string]interface{}
		variant, ok := item.(map[string]interface{})
		if !ok {
			fmt.Println("⚠️ A variant is not a map[string]interface{}")
			continue
		}

		// Safely get the fields and perform necessary type assertions
		quantity, _ := variant["InventoryQuantity"].(float64) // type assertion to float64
		sku, _ := variant["SKU"].(string)
		title, _ := variant["Title"].(string)
		price, _ := variant["Price"].(string)
		bodyHTML, _ := variant["BodyHTML"].(string)

		products = append(products, map[string]interface{}{
			"sku":   sku,
			"name":  title,
			"price": price,
			"status": 1,
			"visibility": 4,
			"type_id": "simple",
			"weight": 1.0,
			"extension_attributes": map[string]interface{}{
				"stock_item": map[string]interface{}{
					"qty":         quantity,
					"is_in_stock": quantity > 0,
				},
			},
			"custom_attributes": []map[string]interface{}{
				{
					"attribute_code": "description",
					"value":          bodyHTML,
				},
			},
		})
	}

	jsonStr, _ := json.MarshalIndent(products, "", "  ")
	fmt.Println("🔥 products:", string(jsonStr))
	fmt.Println("🔥 products:", products)

	// Marshal products to JSON string and return
	return string(jsonStr)
}


func createProduct(product map[string]interface{}) error {
	// Convert the product map into JSON
	productJSON, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("failed to marshal product: %v", err)
	}

	// API URL where the product should be created
	apiURL := "https://api.your-ecommerce-platform.com/products" // Replace with actual API URL

	// Create a new HTTP request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(productJSON))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set necessary headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer YOUR_ACCESS_TOKEN") // Replace with your access token

	// Create an HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second, // 10 seconds timeout
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API call failed with status code: %d", resp.StatusCode)
	}

	// Optionally, you can read the response body if needed
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("Response body:", string(body))

	return nil
}



func main() {
	lambda.Start(HandleProductRequest)
}




