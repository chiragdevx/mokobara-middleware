package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func makeRequest(url, token string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Shopify-Access-Token", token)

	return http.DefaultClient.Do(req)
}

// func fetchResponseBody(resp *http.Response) ([]byte, error) {
// 	defer resp.Body.Close()

// 	if resp.StatusCode != 200 {
// 		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
// 	}

// 	return io.ReadAll(resp.Body)
// }

func getShopifyOrderId(orderID string) string {
	// If the order is found in Shopify, it will return true
	// If not found, it will return false

	storeName := os.Getenv("STORE_NAME")
	shopifyToken := os.Getenv("SHOPIFY_TOKEN")

	// Fetch order
	orderURL := fmt.Sprintf("https://%s.myshopify.com/admin/api/2023-04/orders.json?tag=%s", storeName, orderID)

	// check if the order exists in Shopify
	res, err := makeRequest(orderURL, shopifyToken)

	if err != nil {
		fmt.Printf("❌ Error making request: %v\n", err)
		return ""
	}

	// Check the response status code
	if res.StatusCode != http.StatusOK {
		fmt.Printf("❌ Error: Status code %d\n", res.StatusCode)
		return ""
	}

	// Parse the response
	var responseData struct {
		Orders []struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
			Tags string `json:"tags"`
		} `json:"orders"`
	}
	err = json.NewDecoder(res.Body).Decode(&responseData)
	if err != nil {
		fmt.Printf("❌ Error decoding response: %v\n", err)
		return ""
	}

	// Check if any orders were returned
	if len(responseData.Orders) > 0 {
		fmt.Printf("✅ Order with orderID '%s' found: %s\n", orderID, responseData.Orders[0].Name)
		return fmt.Sprintf("%d", responseData.Orders[0].ID)
	}

	fmt.Printf("❌ No order found with orderID '%s'\n", orderID)
	return ""
}

func getOrderStatus(orderID string) string {
	apiURL := os.Getenv("BASE_URL")
	if apiURL == "" {
		log.Println("updateProduct: ❌ BASE_URL environment variable not set")
		return ""
	}

	req, err := http.NewRequest("GET", apiURL+"/rest/V1/orders/"+orderID+"/status", nil)
	if err != nil {
		log.Printf("❌ Failed to create request: %v\n", err)
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	token := os.Getenv("URL_TOKEN")
	if token == "" {
		log.Println("❌ URL_TOKEN environment variable not set")
		return ""
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Failed to send request: %v\n", err)
		return ""
	}
	defer resp.Body.Close()
	log.Printf("🌐 Response Status: %d\n", resp.StatusCode)

	var statusResponse struct {
		Status      string `json:"status"`
		StatusLabel string `json:"status_label"`
	}

	err = json.NewDecoder(resp.Body).Decode(&statusResponse)
	if err != nil {
		log.Printf("❌ Error decoding response: %v\n", err)
		return ""
	}

	return statusResponse.Status
}
