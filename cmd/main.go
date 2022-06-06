package main

import (
	"fmt"
	rest_api "github.com/fir1/rock/server"
	stock_provider "github.com/fir1/rock/stock-provider"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	logger *logrus.Logger
)

func main() {
	if err := run(); err != nil {
		logger.Error(err)
	}
}

func run() error {
	logger = NewTextLogger()

	stockSymbol, err := getEnvStr("SYMBOL")
	if err != nil {
		return err
	}

	stockNumberOfDays, err := getEnvInt("NDAYS")
	if err != nil {
		return err
	}

	alphaVantageHost, _ := getEnvStr("ALPHA_VANTAGE_HOST")
	if alphaVantageHost == "" {
		alphaVantageHost = "https://www.alphavantage.co"
	}

	alphaVantageApiKey, err := getEnvStr("ALPHA_VANTAGE_API_KEY")
	if err != nil {
		return err
	}

	serverPort, _ := getEnvInt("SERVER_PORT")
	if serverPort == 0 {
		serverPort = 8080
	}

	logger := NewTextLogger()

	stockProvider := stock_provider.NewAlphaVantage(alphaVantageHost, alphaVantageApiKey)

	restApi := rest_api.NewService(logger, serverPort, stockProvider, stockSymbol, stockNumberOfDays)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(interrupt)

	var wg sync.WaitGroup

	wg.Add(1)

	stop := make(chan struct{})
	errChan := make(chan error)

	go func() {
		defer wg.Done()

		err := restApi.Run(stop)
		if err != nil {
			errChan <- err
		}
	}()

	// Wait signal or error from services
	select {
	case <-interrupt:
	case err := <-errChan:
		return fmt.Errorf("webhook rest api server is down (error: %w)", err)
	}

	stop <- struct{}{}
	wg.Wait()

	return nil
}
