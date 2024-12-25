package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Order represents the structure of the incoming order payload.
type Order struct {
	OrderID       string  `json:"order_id"`
	CustomerEmail string  `json:"customer_email"`
	Billing       Address `json:"billing_address"`
	Shipping      Address `json:"shipping_address"`
	Items         []Item  `json:"items"`
}

// Address represents a customer's address.
type Address struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Street    string `json:"street"`
	City      string `json:"city"`
	Region    string `json:"region"`
	Postcode  string `json:"postcode"`
	CountryID string `json:"country_id"`
	Telephone string `json:"telephone"`
}

// Item represents an order item.
type Item struct {
	SKU      string  `json:"sku"`
	Name     string  `json:"name"`
	Quantity int     `json:"qty"`
	Price    float64 `json:"price"`
	RowTotal float64 `json:"row_total"`
}

// ShopifyOrder represents the payload structure for Shopify.
type ShopifyOrder struct {
	Order ShopifyOrderDetails `json:"order"`
}

// ShopifyOrderDetails contains order details for Shopify.
type ShopifyOrderDetails struct {
	Email           string                 `json:"email"`
	Fulfillment     string                 `json:"fulfillment_status"`
	LineItems       []ShopifyLineItem      `json:"line_items"`
	ShippingAddress ShopifyShippingAddress `json:"shipping_address"`
	BillingAddress  ShopifyBillingAddress  `json:"billing_address"`
}

// ShopifyLineItem represents an order item for Shopify.
type ShopifyLineItem struct {
	Title    string `json:"title"`
	Quantity int    `json:"quantity"`
	Price    string `json:"price"`
}

// ShopifyShippingAddress represents the shipping address for Shopify.
type ShopifyShippingAddress struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Address1  string `json:"address1"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	Zip       string `json:"zip"`
	Phone     string `json:"phone"`
}

// ShopifyBillingAddress represents the billing address for Shopify.
type ShopifyBillingAddress struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Address1  string `json:"address1"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	Zip       string `json:"zip"`
	Phone     string `json:"phone"`
}

// HandleOrderRequest handles API Gateway requests
func HandleOrderRequest(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	fmt.Printf("🔥 Received event: %+v\n", request)
	fmt.Printf("🔥 Request body: %s\n", request.Body)

	var order Order
	err := json.Unmarshal([]byte(request.Body), &order)
	if err != nil {
		fmt.Printf("❌ Error unmarshalling body: %v\n", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf(`{"error": "Invalid JSON: %v"}`, err),
		}, nil
	}

	fmt.Printf("🔥 Order: %+v\n", order)

	shopifyOrderId := getShopifyOrderId(order.OrderID)

	if shopifyOrderId == "" {
		createShopifyOrder(order)
	} else {
		updateShopifyOrder(order)
	}

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
