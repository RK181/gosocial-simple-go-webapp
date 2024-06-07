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

// Middleware para manejar las páginas de error
func HandleErrorPage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Crear una copia del ResponseWriter para verificar el código de estado más tarde
			rec := httptest.NewRecorder()
			next.ServeHTTP(rec, r)

			// Verificar el código de estado
			if status := rec.Result().StatusCode; status >= 400 && status <= 599 {
				data := make(map[string]interface{})
				data["StatusCode"] = status
				data["Message"] = http.StatusText(status)
				// Establecer el código de estado en 200
				w.WriteHeader(200)
				// Establecer las cookies
				coockies := rec.Result().Cookies()
				for _, cookie := range coockies {
					http.SetCookie(w, cookie)
				}
				// Escribir la respuesta
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
