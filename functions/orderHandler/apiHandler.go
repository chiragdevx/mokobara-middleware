package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func createShopifyOrder(order Order) error {
	storeName := os.Getenv("STORE_NAME")
	shopifyToken := os.Getenv("SHOPIFY_TOKEN")
	shopifyURL := fmt.Sprintf("https://%s.myshopify.com/admin/api/2021-07/orders.json", storeName)

	shopifyLineItems := []ShopifyLineItem{}

	for _, item := range order.Items {
		shopifyLineItems = append(shopifyLineItems, ShopifyLineItem{
			Title:    item.Name,
			Quantity: item.Quantity,
			Price:    fmt.Sprintf("%.2f", item.Price),
		})
	}

	shopifyOrder := ShopifyOrder{
		Order: ShopifyOrderDetails{
			Email:       order.CustomerEmail,
			Fulfillment: "unfulfilled",
			LineItems:   shopifyLineItems,
			ShippingAddress: ShopifyShippingAddress{
				FirstName: order.Shipping.Firstname,
				LastName:  order.Shipping.Lastname,
				Address1:  order.Shipping.Street,
				City:      order.Shipping.City,
				Province:  order.Shipping.Region,
				Zip:       order.Shipping.Postcode,
				Country:   order.Shipping.CountryID,
			},
			BillingAddress: ShopifyBillingAddress{
				FirstName: order.Billing.Firstname,
				LastName:  order.Billing.Lastname,
				Address1:  order.Billing.Street,
				City:      order.Billing.City,
				Province:  order.Billing.Region,
				Zip:       order.Billing.Postcode,
				Country:   order.Billing.CountryID,
			},
		},
	}

	shopifyOrderJSON, err := json.Marshal(shopifyOrder)
	if err != nil {
		return fmt.Errorf("❌ failed to marshal Shopify order: %v", err)
	}

	// print the order JSON
	fmt.Printf("🔥 Shopify Order JSON: %s\n", shopifyOrderJSON)

	// create a new HTTP request
	req, err := http.NewRequest("POST", shopifyURL, bytes.NewBuffer(shopifyOrderJSON))

	if err != nil {
		return fmt.Errorf("❌ failed to create Shopify request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Shopify-Access-Token", shopifyToken)

	// send the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("❌ failed to send Shopify request: %v", err)
	}

	// check the response status code
	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("❌ unexpected status code: %d", res.StatusCode)
	}

	defer res.Body.Close()

	return nil
}

func updateShopifyOrder(order Order) error {
	storeName := os.Getenv("STORE_NAME")
	shopifyToken := os.Getenv("SHOPIFY_TOKEN")
	shopifyURL := fmt.Sprintf("https://%s.myshopify.com/admin/api/2021-07/orders/%s.json", storeName, order.OrderID)

	status := getOrderStatus(order.OrderID)

	if status == "" {
		return fmt.Errorf("status not found for order: %s", order.OrderID)
	}

	// get the status of the order

	shopifyLineItems := []ShopifyLineItem{}

	for _, item := range order.Items {
		shopifyLineItems = append(shopifyLineItems, ShopifyLineItem{
			Title:    item.Name,
			Quantity: item.Quantity,
			Price:    fmt.Sprintf("%.2f", item.Price),
		})
	}

	shopifyOrder := ShopifyOrder{
		Order: ShopifyOrderDetails{
			Email:       order.CustomerEmail,
			Fulfillment: status,
			LineItems:   shopifyLineItems,
			ShippingAddress: ShopifyShippingAddress{
				FirstName: order.Shipping.Firstname,
				LastName:  order.Shipping.Lastname,
				Address1:  order.Shipping.Street,
				City:      order.Shipping.City,
				Province:  order.Shipping.Postcode,
				Country:   order.Shipping.CountryID,
			},
			BillingAddress: ShopifyBillingAddress{
				FirstName: order.Billing.Firstname,
				LastName:  order.Billing.Lastname,
				Address1:  order.Billing.Street,
				City:      order.Billing.City,
				Province:  order.Billing.Postcode,
				Country:   order.Billing.CountryID,
			},
		},
	}

	shopifyOrderJSON, err := json.Marshal(shopifyOrder)
	if err != nil {
		return fmt.Errorf("❌ failed to marshal Shopify order: %v", err)
	}

	// print the order JSON
	fmt.Printf("🔥 Shopify Order JSON: %s\n", shopifyOrderJSON)

	// create a new HTTP request to update the order
	req, err := http.NewRequest("PUT", shopifyURL, bytes.NewBuffer(shopifyOrderJSON))
	if err != nil {
		return fmt.Errorf("❌ failed to create Shopify request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Shopify-Access-Token", shopifyToken)

	// send the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("❌ failed to send Shopify request: %v", err)
	}

	// check the response status code
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("❌ unexpected status code: %d", res.StatusCode)
	}

	defer res.Body.Close()

	return nil
}
