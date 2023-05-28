package lambda

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/glbter/currency-ex/currency/exchanger"
	"github.com/glbter/currency-ex/pkg/serrors"
	"github.com/glbter/currency-ex/stocks"
	"math/rand"
	"net/http"
	"time"
)

const (
	beforeInQuery = "before"
	afterInQuery  = "after"

	tickerIDInPath   = "ticker_id"
	tickerIDsInQuery = "ticker_ids"
)

type PortfolioHandler struct {
	//db            sqlc.DB
	//portfolioRepo stocks.PortfolioRepository
	userIDExtractor UserIDExtractor

	usecases stocks.PortfolioUsecases

	//tickerRepo stocks.TickerRepository
}

func NewPortfolioHandler(
	//db sqlc.DB,
	usecases stocks.PortfolioUsecases,
	//portfolioRepo stocks.PortfolioRepository,
	//tickerRepo stocks.TickerRepository,
	userIDExtractor UserIDExtractor,
) PortfolioHandler {
	return PortfolioHandler{
		//db:              db,
		usecases: usecases,
		//portfolioRepo:   portfolioRepo,
		//tickerRepo:      tickerRepo,
		userIDExtractor: userIDExtractor,
	}
}

func (h PortfolioHandler) GetCountPortfolio(
	ctx context.Context,
	request events.APIGatewayV2HTTPRequest,
) (events.APIGatewayV2HTTPResponse, error) {
	userID, err := h.userIDExtractor.GetUserID(request)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: serrors.GetHttpCodeFrom(err),
		}, err
	}

	state, err := h.usecases.CountPortfolio(ctx, userID, stocks.ExchangeParams{
		ConverFrom: exchanger.USD,
		ConvertTo:  exchanger.USD,
	})
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: serrors.GetHttpCodeFrom(err),
		}, err
	}

	b, err := json.Marshal(state)
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

type TradeTickerRequest struct {
	Amount   float64                `json:"amount"`
	Action   stocks.PortfolioAction `json:"action"`
	TickerID string                 `json:"ticker_id"`
}

func (h PortfolioHandler) TradeTicker(
	ctx context.Context,
	request events.APIGatewayV2HTTPRequest,
) (events.APIGatewayV2HTTPResponse, error) {
	userID, err := h.userIDExtractor.GetUserID(request)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: serrors.GetHttpCodeFrom(err),
		}, err
	}

	var req TradeTickerRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("unmarshal req: %w", err)
	}

	//tickerID, ok := ExtractFromPath(request, tickerIDInPath)
	//if !ok {
	//	return events.APIGatewayV2HTTPResponse{
	//		StatusCode: http.StatusBadRequest,
	//	}, errors.New("get ticker id")
	//}

	// change to ticker name maybe?

	//amount, err := h.portfolioRepo.CountTickerAmount(ctx, h.db, stocks.CountTickerAmountParams{
	//	UserID:    userID,
	//	TickerIDs: []string{tickerID},
	//})
	//if err != nil {
	//	return events.APIGatewayV2HTTPResponse{
	//		StatusCode: serrors.GetHttpCodeFrom(err),
	//	}, err
	//}
	//
	//if len(amount) != 1 || amount[0].Amount-req.Amount < 0 {
	//	return events.APIGatewayV2HTTPResponse{
	//		StatusCode: http.StatusBadRequest,
	//	}, errors.New("not enough tickers")
	//}
	//
	//dailies, err := h.tickerRepo.QueryLatestDaily(ctx, h.db, stocks.QueryDailyFilter{
	//	TickerIDs: []string{tickerID},
	//})
	//if err != nil {
	//	return events.APIGatewayV2HTTPResponse{
	//		StatusCode: serrors.GetHttpCodeFrom(err),
	//	}, err
	//}
	//
	//if len(dailies) != 0 {
	//	return events.APIGatewayV2HTTPResponse{
	//		StatusCode: serrors.GetHttpCodeFrom(err),
	//	}, err
	//}
	//
	//daily := dailies[0]
	//price := simulatePrice(daily.Low, daily.High)
	//// TODO: add trade by price limit
	//if err := h.portfolioRepo.TradeTickers(ctx, h.db, stocks.TradeTickerParams{
	//	TickerID: tickerID,
	//	UserID:   userID,
	//	Amount:   req.Amount,
	//	PriceUSD: price,
	//	Action:   req.Action,
	//}); err != nil {
	//	return events.APIGatewayV2HTTPResponse{
	//		StatusCode: serrors.GetHttpCodeFrom(err),
	//	}, err
	//}

	if err := h.usecases.TradeTickers(ctx, stocks.TradeTickerParams{
		//TickerID: tickerID,
		TickerID: req.TickerID,
		UserID:   userID,
		Amount:   req.Amount,
		Action:   req.Action,
	}); err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: serrors.GetHttpCodeFrom(err),
		}, err
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusNoContent,
	}, nil
}

func simulatePrice(low, high float64) float64 {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	mean := (high + low) / 2
	stdDev := (high - low) / 6
	return r.NormFloat64()*stdDev + mean
}
