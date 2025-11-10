package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logging logs HTTP requests
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		isHealthCheck := r.URL.Path == "/health"

		// Wrap response writer to capture status code
		rw := newResponseWriter(w)

		// Skip logging health check requests initially
		if !isHealthCheck {
			log.Printf("[%s] %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		}

		next.ServeHTTP(rw, r)

		// Only log health checks if there's an error (non-200 status)
		if isHealthCheck {
			if rw.statusCode != http.StatusOK {
				log.Printf("[%s] %s %s - %d - %v", r.Method, r.RequestURI, r.RemoteAddr, rw.statusCode, time.Since(start))
			}
			return
		}

		log.Printf("Request completed in %v", time.Since(start))
	})
}
