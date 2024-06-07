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
		if r.Method == http.MethodGet {
			// Create a copy of the ResponseWriter to check the status code later
			rec := httptest.NewRecorder()
			next.ServeHTTP(rec, r)

			// Check the status code
			if status := rec.Result().StatusCode; status >= 400 && status <= 599 {
				data := make(map[string]interface{})
				data["StatusCode"] = status
				data["Message"] = http.StatusText(status)

				w.WriteHeader(200)
				shared.ReturnView(w, r, "error.html", data)
			} else {
				if rec.Code == http.StatusSeeOther {
					path, _ := rec.Result().Location()
					http.Redirect(w, r, path.String(), http.StatusSeeOther)
					return
				}
				w.WriteHeader(rec.Result().StatusCode)
				w.Write(rec.Body.Bytes())
				return
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
