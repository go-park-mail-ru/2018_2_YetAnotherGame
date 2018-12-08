package routes

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"2018_2_YetAnotherGame/ApiMS/controllers"
	"2018_2_YetAnotherGame/ApiMS/middlewares"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

func Router(env *controllers.Environment) http.Handler {
	prometheus.MustRegister(env.Counter)
	routerAuth := mux.NewRouter()
	routerAuth.HandleFunc("/api/users/me", env.MeHandle).Methods("GET")
	routerAuth.HandleFunc("/api/users/me", env.UpdateHandle).Methods("POST")
	routerAuth.HandleFunc("/api/upload", env.AvatarHandle).Methods("POST")
	routerAuth.HandleFunc("/api/session", env.LogOutHandle).Methods("DELETE")
	authHandler := middlewares.AuthMiddleware(routerAuth, env.Conn)

	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler())

	router.Handle("/api/users/me", authHandler)
	router.Handle("/api/session", authHandler).Methods("DELETE")
	router.Handle("/api/upload", authHandler)
	router.HandleFunc("/api/leaders", env.ScoreboardHandle).Methods("GET")

	router.HandleFunc("/api/session/new", env.RegistrationHandle).Methods("POST")
	router.HandleFunc("/api/session", env.LoginHandle).Methods("POST")
	router.HandleFunc("/api/vkauth", env.VKRegister)
	return router
}
