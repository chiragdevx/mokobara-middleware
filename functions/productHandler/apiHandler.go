package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func getProductPayload(body map[string]interface{}) ([]map[string]interface{}, error) {
	productData, ok := body["product"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("❌ invalid product data structure")
	}

	variants, ok := productData["variants"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("❌ no variants found in product")
	}
	slug, _ := productData["handle"].(string)

	products := []map[string]interface{}{}
	for _, item := range variants {
		variant, ok := item.(map[string]interface{})
		if !ok {
			fmt.Println("❌ Invalid variant structure")
			continue
		}

		sku, _ := variant["sku"].(string)
		title, _ := variant["title"].(string)
		price, _ := variant["price"].(string)
		inventoryQuantity := float64(0)
		if qty, ok := variant["inventory_quantity"].(float64); ok {
			inventoryQuantity = qty
		}

		combinedSKU := fmt.Sprintf("%s#%s", slug, sku)

		products = append(products, map[string]interface{}{
			"product": map[string]interface{}{
				"sku":              combinedSKU,
				"name":             title,
				"price":            price,
				"status":           1,
				"visibility":       4,
				"type_id":          "simple",
				"weight":           1.0,
				"attribute_set_id": 4,
				"extension_attributes": map[string]interface{}{
					"stock_item": map[string]interface{}{
						"qty":         inventoryQuantity,
						"is_in_stock": inventoryQuantity > 0,
					},
				},
				"custom_attributes": []map[string]interface{}{
					{
						"attribute_code": "description",
						"value":          productData["body_html"],
					},
				},
			},
		})
	}

	return products, nil
}

func createProduct(product map[string]interface{}) error {

	productJSON, err := json.Marshal(product)
	if err != nil {
		log.Printf("❌ Failed to marshal product: %v\n", err)
		return fmt.Errorf("failed to marshal product: %v", err)
	}

	apiURL := os.Getenv("BASE_URL")
	if apiURL == "" {
		log.Println("❌ BASE_URL environment variable not set")
		return fmt.Errorf("BASE_URL environment variable not set")
	}
	log.Printf("🔥 Product JSON: %s\n", productJSON)

	req, err := http.NewRequest("POST", apiURL+"/rest/V1/products", bytes.NewBuffer(productJSON))
	if err != nil {
		log.Printf("❌ Failed to create request: %v\n", err)
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	token := os.Getenv("URL_TOKEN")
	if token == "" {
		log.Println("❌ URL_TOKEN environment variable not set")
		return fmt.Errorf("URL_TOKEN environment variable not set")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Failed to send request: %v\n", err)
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ Failed to read response: %v\n", err)
		return fmt.Errorf("failed to read response: %v", err)
	}

	log.Printf("🌐 Response Status: %d\n", resp.StatusCode)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("❌ API call failed with status code %d: %s\n", resp.StatusCode, string(body))
		return fmt.Errorf("API call failed with status code %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("✅ Product created successfully: %s\n", product["sku"])
	return nil
}
