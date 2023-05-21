package rater

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/glbter/currency-ex/currency/exchanger"
)

type Rater struct {
	c            exchanger.CurrencyRater
	timeInterval time.Duration
}

func (r Rater) FindRates(c exchanger.Currency, start time.Time, end time.Time) ([]exchanger.CurrencyRate, error) {
	var (
		wg sync.WaitGroup

		intervals = r.countAmountOfIntervals(start, end, r.timeInterval)

		rates = make(chan exchanger.CurrencyRate, intervals)
		errs  = make(chan error, intervals)
	)

	for i := int64(0); i < intervals; i++ {
		date := start.Add(time.Duration(i) * r.timeInterval)

		wg.Add(1)
		go func(date time.Time) {
			defer wg.Done()

			rate, err := r.c.FindRate(c, date)
			if err != nil {
				errs <- fmt.Errorf("find rate for %v: %w", date, err)
				return
			}

			rates <- rate
		}(date)
	}

	wg.Wait()

	close(rates)
	close(errs)

	res := make([]exchanger.CurrencyRate, 0, intervals)
	for rate := range rates {
		res = append(res, rate)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Date.Before(res[j].Date)
	})

	var err []error
	for e := range errs {
		err = append(err, e)
	}

	if len(err) != 0 {
		return res, errors.Join(err...)
	}

	return res, nil
}

func (r Rater) countAmountOfIntervals(start time.Time, end time.Time, interval time.Duration) int64 {
	return end.Sub(start).Nanoseconds() / interval.Nanoseconds()
}
