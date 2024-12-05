package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context) (string, error) {
	fmt.Printf("🔥  Hello from Go Lambda!")
    return "Hello from Go Lambda!", nil
}

func main() {
    lambda.Start(HandleRequest)
}
