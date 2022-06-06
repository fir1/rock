package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Service) Run(stop chan struct{}) error {
	s.router = mux.NewRouter()
	// Register all routes on http handler
	s.routes()
	handler := cors.Default().Handler(s.router)
	// Create an instance of our LoggingMiddleware with our configured logger
	loggingMiddleware := s.loggingMiddleware()

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", s.serverPort),
		Handler: loggingMiddleware(handler),
	}

	// channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.logger.Printf("REST API listening on port %d", s.serverPort)
		serverErrors <- server.ListenAndServe()
	}()

	// blocking run and waiting for shutdown.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("error: starting REST API server: %w", err)
	case <-stop:
		s.logger.Warn("webhook server receive STOP signal")
		// asking listener to shutdown
		err := server.Shutdown(context.Background())
		if err != nil {
			return fmt.Errorf("graceful shutdown did not complete: %w", err)
		}
		s.logger.Info("server was shut down gracefully")
	}
	return nil
}
