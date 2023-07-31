package middlewares

import (
	"net/http"

	"github.com/joaopegoraro/ahpsico-go/server"
)

func Security(s *server.Server) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(
				"Content-Security-Policy",
				"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
			)
			w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "deny")
			w.Header().Set("X-XSS-Protection", "0")
			w.Header().Set("Access-Control-Allow-Methods", "Allow")
			w.Header().Set("Access-Control-Allow-Origin", "*")

			next.ServeHTTP(w, r)
		})
	}
}
