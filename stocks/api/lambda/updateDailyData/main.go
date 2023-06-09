package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	cpolygon "github.com/glbter/currency-ex/polygon"
	lambdas "github.com/glbter/currency-ex/stocks/api/lambda"
	"time"
)

type Request struct {
	Multiplier int       `json:"multiplier"`
	Timespan   string    `json:"timespan"`
	From       time.Time `json:"from"`
	To         time.Time `json:"to"`
	Tickers    []string  `json:"tickers"`
}

func main() {
	lambda.Start(func(
		ctx context.Context,
		request Request,
	) error {
		processor, err := lambdas.InitLambdaDailyProcessor(ctx)
		if err != nil {
			return err
		}

		//log.Println("BODY: ", request.Body)
		//var req Request
		//if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		//	return fmt.Errorf("unmarshal req: %w", err)
		//}

		req := request

		processor = processor.WithNewConfig(cpolygon.Config{
			Multiplier: req.Multiplier,
			Timespan:   req.Timespan,
			From:       req.From,
			To:         req.To,
		})

		ts := req.Tickers
		if len(ts) > 5 {
			ts = ts[0:5]
		}

		if err := processor.Process(ctx, ts); err != nil {
			return err
		}

		return nil
	})
}
