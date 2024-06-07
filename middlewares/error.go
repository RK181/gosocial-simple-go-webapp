package middlewares

import (
	"module/shared"
	"net/http"
	"net/http/httptest"
)

type ErrorData struct {
	StatusCode int
	Message    string
}

// Middleware to handle the error pages
func HandleErrorPage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a copy of the ResponseWriter to check the status code later
		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)

		// Check the status code
		if status := rec.Code; status >= 400 && status <= 599 {
			data := make(map[string]interface{})
			data["StatusCode"] = status
			data["Message"] = http.StatusText(status)
			w.WriteHeader(200)

			shared.ReturnView(w, r, "error.html", data)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
