package middlewares

import (
	"context"
	"module/models"
	"module/shared"
	"net/http"
)

// Función que se encarga de redirigir al usuario a la página de login
func writeUnauthed(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Función que se encarga de redirigir al usuario a la página de inicio
func redirectFromBaseToHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// Middleware que se encarga de recuperar y adjuntar la información de autenticación
func FetchAuthInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			redirectFromBaseToHome(w, r)
		}

		isLoginORRegister := r.URL.Path == "/login" || r.URL.Path == "/register"
		cookie, err := r.Cookie(shared.AUTH_USER_TOKEN)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		userSessionToken := cookie.Value
		user := models.User{}
		if !user.GetUserBySessionToken(userSessionToken) {
			next.ServeHTTP(w, r)
			return
		}

		if isLoginORRegister {
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}

		//fmt.Println("User Token: ", userSessionToken)
		ctx := context.WithValue(r.Context(), shared.AUTH_USER, user)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

// Middleware que se encarga de comprobar si el usuario esta autenticado y restrinje el acceso a las rutas
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(shared.AUTH_USER) == nil {
			writeUnauthed(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
