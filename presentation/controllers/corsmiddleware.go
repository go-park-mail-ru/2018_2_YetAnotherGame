package controllers

import (
	"net/http"

	"github.com/rs/cors"
)

func CORSMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := cors.New(cors.Options{
			AllowCredentials: true,
			AllowedOrigins:   []string{"http://127.0.0.1:3000"},
			AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"},
		}).Handler(handler)
		c.ServeHTTP(w, r)
	})
}
