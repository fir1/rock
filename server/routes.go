package server

func (s *Service) routes() {
	s.router.HandleFunc("/healthz", s.handleHealth())
	s.router.HandleFunc("/stock", s.handleGetStock()).Methods("GET")
}
