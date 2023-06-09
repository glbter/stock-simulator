package postgres

import (
	"context"
	"fmt"
	sqlc "github.com/glbter/currency-ex/pkg/sql"
	"github.com/glbter/currency-ex/stocks"
	"github.com/huandu/go-sqlbuilder"
	"strings"
)

func NewPortfolioRepository() PortfolioRepository {
	return PortfolioRepository{}
}

var _ stocks.PortfolioRepository = PortfolioRepository{}

type PortfolioRepository struct{}

type portfolio struct {
	TickerID    string  `db:"ticker_id"`
	Amount      float64 `db:"amount"`
	Name        string  `db:"name"`
	Description string  `db:"description"`
	High        float64 `db:"high"`
	Low         float64 `db:"low"`
	Open        float64 `db:"open"`
	Close       float64 `db:"close"`
}

func (p portfolio) toEntity() stocks.PortfolioTickerState {
	return stocks.PortfolioTickerState{
		TickerID:    p.TickerID,
		Amount:      p.Amount,
		Name:        p.Name,
		Description: p.Description,
		High:        p.High,
		Low:         p.Low,
		Open:        p.Open,
		Close:       p.Close,
	}
}

type tickerAmount struct {
	TickerID string  `db:"ticker_id"`
	Amount   float64 `db:"amount"`
}

func (p tickerAmount) toEntity() stocks.PortfolioTickerAmount {
	return stocks.PortfolioTickerAmount{
		TickerID: p.TickerID,
		Amount:   p.Amount,
	}
}

func (PortfolioRepository) TradeTickers(ctx context.Context, e sqlc.Executor, p stocks.TradeTickerParams) error {
	ib := sqlbuilder.NewInsertBuilder()
	ib.InsertInto("portfolio_record").
		Cols(
			"investor_id",
			"ticker_id",
			"amount",
			"price_usd",
			"action",
		).
		Values(
			p.UserID,
			p.TickerID,
			p.Amount,
			p.PriceUSD,
			p.Action,
		)

	q, args := ib.BuildWithFlavor(sqlbuilder.PostgreSQL)

	if _, err := e.Exec(ctx, q, args...); err != nil {
		return err
	}

	return nil
}

func (PortfolioRepository) CountTickerAmount(
	ctx context.Context,
	s sqlc.Selector,
	p stocks.CountTickerAmountParams,
) ([]stocks.PortfolioTickerAmount, error) {
	tickerSQLFilter := ""
	if len(p.TickerIDs) > 0 {
		l := make([]string, 0, len(p.TickerIDs))
		for i := 0; i < len(p.TickerIDs); i++ {
			l = append(l, fmt.Sprintf("$%d", i+2))
		}
		tickerSQLFilter = fmt.Sprintf("and pr.ticker_id IN (%v)", strings.Join(l, ", "))
	}

	q := `
select
    t.ticker_id,
    coalesce(t.total_amount, 0) as amount
from ticker tt
join (
	select
		t.ticker_id,
		sum(t.mul * pr.amount) as total_amount
	from portfolio_record as pr
	left join (
		-- count for each record the multiply value
		select
			pr.id as check_id,
			pr.ticker_id,
			(case
				 when pr.action = 'BUY' then mul(coalesce(s.after / s.before, 1))
				 when pr.action = 'SELL' then -mul(coalesce(s.after / s.before, 1))
			end) as mul
		from split as s
		right join portfolio_record as pr
			on s.ticker_id = pr.ticker_id
			and s.date > pr.date
		where pr.investor_id = $1 ` + tickerSQLFilter +
		`group by pr.id, pr.action, pr.ticker_id
	) as t
	  on pr.id = t.check_id
	group by t.ticker_id
) as t
  on t.ticker_id = tt.id
`
	args := make([]any, 0, len(p.TickerIDs)+1)
	args = append(args, p.UserID)
	for _, t := range p.TickerIDs {
		args = append(args, t)
	}

	var sqlRes []tickerAmount
	if err := s.Select(ctx, &sqlRes, q, args...); err != nil {
		return nil, err
	}

	res := make([]stocks.PortfolioTickerAmount, 0, len(sqlRes))
	for _, s := range sqlRes {
		res = append(res, s.toEntity())
	}

	return res, nil
}

func (PortfolioRepository) CountPortfolio(ctx context.Context, s sqlc.Selector, userID string) (stocks.PortfolioState, error) {
	q := `
-- collect portfolio state (current amount of all tickers) and current net worth
with portfolio_state as (
    select
        t.ticker_id,
        t.total_amount               as amount,
        tt.name,
        tt.description,
        t.total_amount * stock.high  as high,
        t.total_amount * stock.low   as low,
        t.total_amount * stock.open  as open,
        t.total_amount * stock.close as close
    from ticker tt
	join (
        select stock_daily.ticker_id, high, low, open, close
		from stock_daily
		join (
			select ticker_id, max(date) as date
			from stock_daily
			group by ticker_id
		) as latest_daily
			on stock_daily.ticker_id = latest_daily.ticker_id
			and stock_daily.date = latest_daily.date
    ) stock
	  on tt.id = stock.ticker_id
	join (
        select
            t.ticker_id,
            sum(t.mul * pr.amount) as total_amount
        from portfolio_record as pr
		join (
            -- count for each record the multiply value
            select
                pr.id as check_id,
                pr.ticker_id,
                (case
                     when pr.action = 'BUY' then mul(coalesce(s.after / s.before, 1))
                     when pr.action = 'SELL' then -mul(coalesce(s.after / s.before, 1))
				end) as mul
            from split as s
			right join portfolio_record as pr
				on s.ticker_id = pr.ticker_id
				and s.date > pr.date
            where pr.investor_id = $1
            group by pr.id, pr.action, pr.ticker_id
        ) as t
		  on pr.id = t.check_id
        group by t.ticker_id
    ) as t
	  on t.ticker_id = tt.id
)

select
    ticker_id,
    coalesce(amount, 0) 		as amount,
    name,
    coalesce(description, '') 	as description,
    coalesce(high, 0)			as high,
    coalesce(low, 0) 			as low, 
    coalesce(open, 0) 			as open,
    coalesce(close, 0) 			as close
from (
	 select
		 ticker_id,
		 amount,
		 name,
		 description,
		 high,
		 low,
		 open,
		 close,
		 1 as ordering
	 from portfolio_state
	 union all
	 select
		 uuid_nil()  as ticker_id,
		 sum(amount) as amount,
		 'TOTAL'     as name,
		 'TOTAL'     as description,
		 sum(high)   as high,
		 sum(low)    as low,
		 sum(open)   as open,
		 sum(close)  as close,
		 2           as ordering
	 from portfolio_state
	 order by ordering, name
) t;
`
	var sqlRes []portfolio
	if err := s.Select(ctx, &sqlRes, q, userID); err != nil {
		return stocks.PortfolioState{}, err
	}

	res := stocks.PortfolioState{
		All:   make([]stocks.PortfolioTickerState, 0, len(sqlRes)-1),
		Total: sqlRes[len(sqlRes)-1].toEntity(),
	}

	for _, s := range sqlRes[:len(sqlRes)-1] {
		res.All = append(res.All, s.toEntity())
	}

	return res, nil
}
