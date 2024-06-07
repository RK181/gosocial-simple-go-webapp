package middlewares

import (
	"log"
	"net/http"
	"time"
)

type loggingWrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *loggingWrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

// Middleware to log the request and response
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lww := &loggingWrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(lww, r)
		log.Printf("Status: %d | Method: %s | Path: %s | Consumed Time: %v \n", lww.statusCode, r.Method, r.URL.Path, time.Since(start))

	})
}
