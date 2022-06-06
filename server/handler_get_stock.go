package server

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
)

func (s *Service) handleGetStock() http.HandlerFunc {
	type getStockRequest struct {
		Symbol string `json:"symbol"`
		Ndays  int    `json:"ndays"`
	}

	type dailyPrice struct {
		Open   string `json:"open"`
		High   string `json:"high"`
		Low    string `json:"low"`
		Close  string `json:"close"`
		Volume string `json:"volume"`
	}

	type getStockResponse struct {
		MetaData            map[string]interface{} `json:"meta_data"`
		TimeSeriesDaily     map[string]dailyPrice  `json:"time_series_daily"`
		AverageClosingPrice string                 `json:"average_closing_price"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		request := getStockRequest{}
		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := s.decode(r, &request)
		// in case if request is empty it will return 'EOF' error on that case
		// use environment variable data instead of taking from request
		if err != nil && err.Error() != "EOF" {
			s.respond(w, err.Error(), http.StatusBadRequest)
			return
		}

		symbol := request.Symbol
		if symbol == "" {
			symbol = s.stockSymbol
		}

		ndays := request.Ndays
		if ndays == 0 {
			ndays = s.stockNumberOfDays
		}

		ctx := r.Context()
		getStockPriceResponse, err := s.stockProvider.GetStockPrice(ctx, symbol)
		if err != nil {
			s.respond(w, nil, http.StatusInternalServerError)
			return
		}

		var totalClosingPrice float64
		var count int

		timeSeriesDailyPrice := make(map[string]dailyPrice, ndays)

		// by default when unmarshalling key, value to the map the order of it will be lost, golang does not support unmarshalling with order
		// so here we will make manipulation to sort key and value
		keys := make([]string, 0, len(getStockPriceResponse.TimeSeriesDaily))
		for k := range getStockPriceResponse.TimeSeriesDaily {
			keys = append(keys, k)
		}
		// sort by increasing order date from oldest to the newest
		sort.Strings(keys)

		for i := len(getStockPriceResponse.TimeSeriesDaily) - 1; i >= 0; i-- {
			if count == ndays {
				break
			}

			key := keys[i]
			price := getStockPriceResponse.TimeSeriesDaily[key]

			timeSeriesDailyPrice[key] = dailyPrice{
				Open:   price.Open,
				High:   price.High,
				Low:    price.Low,
				Close:  price.Close,
				Volume: price.Volume,
			}

			closingPrice, err := strconv.ParseFloat(price.Close, 64)
			if err != nil {
				s.respond(w, nil, http.StatusInternalServerError)
				return
			}

			totalClosingPrice += closingPrice
			count++
		}

		averageClosingPrice := totalClosingPrice / float64(count)

		s.respond(w, getStockResponse{
			MetaData:            getStockPriceResponse.MetaData,
			TimeSeriesDaily:     timeSeriesDailyPrice,
			AverageClosingPrice: fmt.Sprintf("%.4f", averageClosingPrice),
		}, http.StatusOK)
	}
}
