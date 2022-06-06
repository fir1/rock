package server

import (
	"fmt"
	"net/http"
	"time"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging. This type will implement http.ResponseWriter.
type responseWriter struct {
	http.ResponseWriter
	status      int
	body        []byte
	wroteHeader bool
	wroteBody   bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteBody {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func (rw *responseWriter) Write(body []byte) (int, error) {
	if rw.wroteBody {
		return 0, nil
	}
	i, err := rw.ResponseWriter.Write(body)
	if err != nil {
		return 0, err
	}
	rw.body = body
	return i, err
}

func (rw *responseWriter) Body() []byte {
	return rw.body
}

// LoggingInDetailMiddleware logs the incoming HTTP request and response. Enable it only for debug purpose disable it on production.
func (s *Service) loggingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/healthz" {
				// Call the next handler don't log if it is internal request from health check of Kubernetes
				next.ServeHTTP(w, r)
				return
			}

			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			requestBody, err := s.readRequestBody(r)
			if err != nil {
				s.respond(w, err.Error(), http.StatusInternalServerError)
				return
			}
			s.restoreRequestBody(r, requestBody)
			logMessage := fmt.Sprintf("path:%s, method: %s", r.URL.EscapedPath(), r.Method)
			requestBodyInStr := string(requestBody)
			if requestBodyInStr == "" {
				logMessage = fmt.Sprintf("%s, requestBody: empty", logMessage)
			} else {
				logMessage = fmt.Sprintf("%s, requestBody: %v", logMessage, requestBodyInStr)
			}

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			s.logger.Infof("request %s, duration: %v", logMessage, time.Since(start))
			s.logger.Infof("response status_code=%d and body=%s", wrapped.Status(), string(wrapped.Body()))
		}
		return http.HandlerFunc(fn)
	}
}
