package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/glbter/currency-ex/pkg/serrors"
	sqlc "github.com/glbter/currency-ex/sql"
	"github.com/glbter/currency-ex/stocks"
	"net/http"
	"time"
)

type TickerHandler struct {
	db sqlc.DB

	tickerRepo stocks.TickerRepository
}

func NewTickerHandler(
	db sqlc.DB,
	tickerRepo stocks.TickerRepository,
) TickerHandler {
	return TickerHandler{
		db:         db,
		tickerRepo: tickerRepo,
	}
}

func (h TickerHandler) GetTickers(
	ctx context.Context,
	request events.APIGatewayV2HTTPRequest,
) (events.APIGatewayV2HTTPResponse, error) {
	tickerIDs, ok := ExtractSliceFromQuery(request, tickerIDsInQuery)
	if !ok {
		tickerIDs = nil
	}

	tickerDaily, err := h.tickerRepo.QueryLatestDaily(ctx, h.db, stocks.QueryDailyFilter{TickerIDs: tickerIDs})
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: serrors.GetHttpCodeFrom(err),
		}, err
	}

	b, err := json.Marshal(tickerDaily)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	return events.APIGatewayV2HTTPResponse{
		Body:       string(b),
		StatusCode: http.StatusOK,
	}, nil
}

func (h TickerHandler) GetTickerGraph(
	ctx context.Context,
	request events.APIGatewayV2HTTPRequest,
) (events.APIGatewayV2HTTPResponse, error) {
	tickerIDs, ok := ExtractSliceFromQuery(request, tickerIDsInQuery)
	if !ok {
		tickerIDs = nil
	}

	params := stocks.QueryDailyGraphParams{
		TickerIDs: tickerIDs,
	}

	beforeS, ok := ExtractFromQuery(request, beforeInQuery)
	if ok {
		before, err := time.Parse(time.DateOnly, beforeS)
		if err != nil {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusBadRequest,
			}, fmt.Errorf("parse date before %w", err)
		}

		params.BeforeDateInc = &before
	}

	afterS, ok := ExtractFromQuery(request, afterInQuery)
	if ok {
		after, err := time.Parse(time.DateOnly, afterS)
		if err != nil {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: http.StatusBadRequest,
			}, fmt.Errorf("parse date after %w", err)
		}

		params.AfterDateInc = &after
	}

	graph, err := h.tickerRepo.QueryTickerDailyGraph(ctx, h.db, params)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: serrors.GetHttpCodeFrom(err),
		}, err
	}

	b, err := json.Marshal(graph)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	return events.APIGatewayV2HTTPResponse{
		Body:       string(b),
		StatusCode: http.StatusOK,
	}, nil
}
