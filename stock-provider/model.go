package stock_provider

import (
	"encoding/json"
)

func UnmarshalGetStockPriceResponse(data []byte) (GetStockPriceResponse, error) {
	var r GetStockPriceResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GetStockPriceResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type GetStockPriceResponse struct {
	MetaData        map[string]interface{}     `json:"Meta Data"`
	TimeSeriesDaily map[string]TimeSeriesDaily `json:"Time Series (Daily)"`
}

type MetaData struct {
	The1Information   string `json:"1. Information"`
	The2Symbol        string `json:"2. Symbol"`
	The3LastRefreshed string `json:"3. Last Refreshed"`
	The4OutputSize    string `json:"4. Output Size"`
	The5TimeZone      string `json:"5. Time Zone"`
}

type TimeSeriesDaily struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}
