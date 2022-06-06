package server

import "net/http"

func (s *Service) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			s.logger.Errorf("health write error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
