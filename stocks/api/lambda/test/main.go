package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
)

type Result struct {
	InsideRadius bool `json:"insideRadius"`
}

func HandleRequest(ctx context.Context) (Result, error) {
	result := Result{InsideRadius: true}
	return result, nil
}

func main() {
	lambda.Start(HandleRequest)
}
