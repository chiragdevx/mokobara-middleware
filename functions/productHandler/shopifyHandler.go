package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

func fetchResponseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func parseMetafields(responseBody []byte) (bool, error) {
	var metafieldResponse map[string]interface{}
	if err := json.Unmarshal(responseBody, &metafieldResponse); err != nil {
		return false, fmt.Errorf("error unmarshalling response: %w", err)
	}

	if metafields, ok := metafieldResponse["metafields"].([]interface{}); ok {
		for _, m := range metafields {
			if mf, ok := m.(map[string]interface{}); ok && mf["key"] == "is_published" {
				if value, ok := mf["value"].(bool); ok && value {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func getProductWithMetafields(productID string) (map[string]interface{}, error) {
	storeName := os.Getenv("STORE_NAME")
	shopifyToken := os.Getenv("SHOPIFY_TOKEN")

	// Fetch metafields
	metafieldURL := fmt.Sprintf("https://%s.myshopify.com/admin/api/2023-04/products/%s/metafields.json", storeName, productID)
	resp, err := makeRequest(metafieldURL, shopifyToken)
	if err != nil {
		return nil, err
	}

	body, err := fetchResponseBody(resp)
	if err != nil {
		return nil, err
	}

	isPublished, err := parseMetafields(body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("🔥 isPublished: %v\n", isPublished)

	if !isPublished {
		return nil, fmt.Errorf("❌ product is not published")
	}

	// Fetch product details if published
	productURL := fmt.Sprintf("https://%s.myshopify.com/admin/api/2023-04/products/%s.json", storeName, productID)

	resp, err = makeRequest(productURL, shopifyToken)
	if err != nil {
		return nil, err
	}

	body, err = fetchResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var productResponse map[string]interface{}
	if err := json.Unmarshal(body, &productResponse); err != nil {
		return nil, fmt.Errorf("❌ error unmarshalling product response: %w", err)
	}
	return productResponse, nil
}
