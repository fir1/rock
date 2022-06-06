package stock_provider

import (
	"context"
)

type ProviderInterface interface {
	GetStockPrice(ctx context.Context, symbolName string) (GetStockPriceResponse, error)
}
