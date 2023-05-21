package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/glbter/currency-ex/pkg/serrors"
	lambdas "github.com/glbter/currency-ex/stocks/api/lambda"
)

func main() {
	lambda.Start(func(
		ctx context.Context,
		request events.APIGatewayV2HTTPRequest,
	) (events.APIGatewayV2HTTPResponse, error) {
		handler, err := lambdas.InitLambdaPortfolioHandler(ctx)
		if err != nil {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: serrors.GetHttpCodeFrom(err),
			}, err
		}

		return handler.GetCountPortfolio(ctx, request)
	})
}
