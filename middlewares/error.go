package middlewares

import (
	"log"
	"module/shared"
	"net/http"
	"net/http/httptest"
)

type ErrorData struct {
	StatusCode int
	Message    string
}

func CatchError(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a copy of the ResponseWriter to check the status code later
		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)

		// Check the status code
		if status := rec.Code; status >= 400 && status <= 599 {
			log.Println("Error 1: ", rec.Code)

		} else {
			log.Println("Error 2: ", rec.Code)
			next.ServeHTTP(w, r)
		}
	})
}

func handleError(w http.ResponseWriter, statusCode int) {
	switch statusCode {
	case 500 - 599:
		tmpl := shared.Templates["error.html"]
		data := ErrorData{
			StatusCode: statusCode,
			Message:    "Internal Server Error",
		}
		tmpl.Execute(w, data)
		return
	case 404:
		tmpl := shared.Templates["error.html"]
		data := ErrorData{
			StatusCode: statusCode,
			Message:    http.StatusText(statusCode),
		}
		w.WriteHeader(200)

		tmpl.Execute(w, data)
		log.Println(http.StatusText(statusCode))
		return
	default:
	}
}
