package routes

import (
	"net/http"

	"2018_2_YetAnotherGame/presentation/controllers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Router(env *controllers.Environment) http.Handler {
	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"http://127.0.0.1:3000"},                           // All origins
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"}, // Allowing only get, just an example

	})
	router := mux.NewRouter()
	router.HandleFunc("/api/session", env.LoginHandle).Methods("POST")
	router.HandleFunc("/api/session/new", env.RegistrationHandle).Methods("POST")
	router.HandleFunc("/api/users/me", env.MeHandle).Methods("GET")
	router.HandleFunc("/api/leaders", env.ScoreboardHandle).Methods("GET")
	router.HandleFunc("/api/session", env.LogOutHandle).Methods("DELETE")
	router.HandleFunc("/api/avatar", env.AvatarHandle)
	return c.Handler(router)
}
