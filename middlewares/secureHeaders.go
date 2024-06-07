package middlewares

import "net/http"

// Middleware to add secure headers to the response
func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubdomains")
		w.Header().Add("Content-Security-Policy", "default-src 'self'; script-src https://cdn.jsdelivr.net; style-src https://cdn.jsdelivr.net; img-src 'self' data:;")
		w.Header().Add("X-XSS-Protection", "1; mode=block")
		w.Header().Add("X-Frame-Options", "DENY")
		w.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Add("X-Content-Type-Options", "nosniff")
		w.Header().Add("Content-Type", "text/html; charset=UTF-8")
		w.Header().Add("Access-Control-Request-Method", "GET, POST")

		next.ServeHTTP(w, r)
	})
}
