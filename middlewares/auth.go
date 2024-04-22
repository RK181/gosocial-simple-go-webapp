package middlewares

import (
	"context"
	"fmt"
	"module/models"
	"net/http"
)

type contextKey string

const (
	AUTH_USER       contextKey = "authUser"
	AUTH_USER_TOKEN string     = "authUserToken"
)

func writeUnauthed(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
func redirectFromBaseToHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isLoginORRegister := r.URL.Path == "/login" || r.URL.Path == "/register"
		basePath := r.URL.Path == "/"
		cookie, err := r.Cookie(AUTH_USER_TOKEN)
		if err != nil {
			if isLoginORRegister || basePath {
				next.ServeHTTP(w, r)
				return
			}
			writeUnauthed(w, r)
			return
		}

		userSessionToken := cookie.Value
		user := models.User{}
		if !user.GetUserBySessionToken(userSessionToken) {
			if isLoginORRegister || basePath {
				next.ServeHTTP(w, r)
				return
			}
			writeUnauthed(w, r)
			return
		}

		if isLoginORRegister {
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}

		fmt.Println("User Token: ", userSessionToken)
		// ctx
		ctx := context.WithValue(r.Context(), AUTH_USER, user)

		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

func FetchAuthInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			redirectFromBaseToHome(w, r)
		}

		isLoginORRegister := r.URL.Path == "/login" || r.URL.Path == "/register"
		cookie, err := r.Cookie(AUTH_USER_TOKEN)
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

		fmt.Println("User Token: ", userSessionToken)
		ctx := context.WithValue(r.Context(), AUTH_USER, user)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(AUTH_USER) == nil {
			writeUnauthed(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
