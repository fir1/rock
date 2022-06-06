package stock_provider

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type AlphaVantage struct {
	host   string
	apiKey string
}

func NewAlphaVantage(host, apiKey string) ProviderInterface {
	return &AlphaVantage{
		host:   host,
		apiKey: apiKey,
	}
}

var httpClient = http.Client{
	Timeout: 1 * time.Minute,
}

func (a AlphaVantage) GetStockPrice(ctx context.Context, symbolName string) (GetStockPriceResponse, error) {
	url := fmt.Sprintf("%s/query?apikey=%s&function=TIME_SERIES_DAILY&symbol=%s", a.host, a.apiKey, symbolName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return GetStockPriceResponse{}, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return GetStockPriceResponse{}, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return GetStockPriceResponse{}, err
	}
	defer resp.Body.Close()

	return UnmarshalGetStockPriceResponse(content)
}
