package stocks

import (
	"fmt"
	"github.com/glbter/currency-ex/pkg/serrors"
	"time"
)

type Ticker struct {
	ID          string
	Name        string
	Description string
	//Type string
	//Market string
}

type TickerWithData struct {
	Ticker Ticker

	High   float64
	Low    float64
	Open   float64
	Close  float64
	Volume float64
}

type TickerSplit struct {
	TickerID string
	Name     string
	Date     time.Time
	From     float64
	To       float64
}

type StockDailyData struct {
	TickerID string
	Date     time.Time
	High     float64
	Low      float64
	Open     float64
	Close    float64
	Volume   float64
}

type PortfolioAction string

const (
	ACTION_BUY  PortfolioAction = "BUY"
	ACTION_SELL PortfolioAction = "SELL"
)

func (a PortfolioAction) Check() error {
	if a != ACTION_BUY && a != ACTION_SELL {
		return fmt.Errorf("no such action: %w", serrors.ErrBadInput)
	}

	return nil
}

type PortfolioRecord struct {
	ID         string
	investorID string
	tickerID   string
	date       time.Time
	price      float64
	action     PortfolioAction
}

type PortfolioTickerAmount struct {
	TickerID string
	Amount   float64
}

type PortfolioState struct {
	Total PortfolioTickerState
	All   []PortfolioTickerState
}

type PortfolioTickerState struct {
	TickerID    string
	Amount      float64
	Name        string
	Description string
	High        float64
	Low         float64
	Open        float64
	Close       float64
}
