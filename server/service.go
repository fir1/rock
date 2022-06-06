package server

import (
	stock_provider "github.com/fir1/rock/stock-provider"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Service struct {
	router            *mux.Router
	logger            *logrus.Logger
	serverPort        int
	stockProvider     stock_provider.ProviderInterface
	stockSymbol       string
	stockNumberOfDays int
}

func NewService(logger *logrus.Logger,
	serverPort int, stp stock_provider.ProviderInterface,
	symbol string, ndays int) *Service {
	return &Service{
		logger:            logger,
		serverPort:        serverPort,
		stockProvider:     stp,
		stockSymbol:       symbol,
		stockNumberOfDays: ndays,
	}
}
