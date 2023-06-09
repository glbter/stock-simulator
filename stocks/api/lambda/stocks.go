package lambda

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/glbter/currency-ex/currency/exchanger"
	"github.com/glbter/currency-ex/pkg/serrors"
	"github.com/glbter/currency-ex/stocks"
	"github.com/glbter/currency-ex/stocks/usecases"
	"net/http"
)

const (
	beforeInQuery = "before"
	afterInQuery  = "after"

	tickerIDInPath   = "ticker_id"
	tickerIDsInQuery = "ticker_ids"
)

type PortfolioHandler struct {
	userIDExtractor UserIDExtractor

	usecases stocks.PortfolioUsecases
}

func NewPortfolioHandler(
	usecases stocks.PortfolioUsecases,
	userIDExtractor UserIDExtractor,
) PortfolioHandler {
	return PortfolioHandler{
		usecases:        usecases,
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
		}, serrors.GetErrorByTypeAndLog(err)
	}

	state, err := h.usecases.CountPortfolio(ctx, userID, stocks.ExchangeParams{
		ConvertFrom: exchanger.USD,
		ConvertTo:   exchanger.UAH,
	})
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: serrors.GetHttpCodeFrom(err),
		}, serrors.GetErrorByTypeAndLog(err)
	}

	b, err := json.Marshal(state)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
		}, serrors.GetErrorByTypeAndLog(err)
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
		}, serrors.GetErrorByTypeAndLog(fmt.Errorf("%w: unmarshal req: %v", serrors.ErrBadInput, err))
	}

	//// TODO: add trade by price limit
	if err := h.usecases.TradeTickers(ctx, stocks.TradeTickerParams{
		//TickerID: tickerID,
		TickerID: req.TickerID,
		UserID:   userID,
		Amount:   req.Amount,
		Action:   req.Action,
	}); err != nil {
		var e usecases.ErrNotEnoughTickers
		if errors.As(err, &e) {
			return events.APIGatewayV2HTTPResponse{
				Body:       fmt.Sprintf("{\"msg\":\"%v\"}", e.Error()),
				StatusCode: http.StatusBadRequest,
			}, nil
		}

		return events.APIGatewayV2HTTPResponse{
			StatusCode: serrors.GetHttpCodeFrom(err),
		}, serrors.GetErrorByTypeAndLog(err)
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
	}, nil
}
