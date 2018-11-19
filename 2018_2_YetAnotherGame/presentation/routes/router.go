package routes

import (
	"net/http"

	"2018_2_YetAnotherGame/presentation/controllers"
	"2018_2_YetAnotherGame/presentation/middlewares"

	"github.com/gorilla/mux"
)

func Router(env *controllers.Environment) http.Handler {
	routerAuth := mux.NewRouter()
	routerAuth.HandleFunc("/api/user/me", env.MeHandle).Methods("GET")
	routerAuth.HandleFunc("/api/user/me", env.UpdateHandle).Methods("POST")
	routerAuth.HandleFunc("/api/avatar", env.AvatarHandle).Methods("POST")
	routerAuth.HandleFunc("/api/session", env.LogOutHandle).Methods("DELETE")
	authHandler := middlewares.AuthMiddleware(routerAuth, env.DB)

	router := mux.NewRouter()
	router.Handle("/api/user/me", authHandler)
	router.Handle("/api/session", authHandler).Methods("DELETE")
	router.Handle("/api/upload", authHandler)
	router.HandleFunc("/api/leaders", env.ScoreboardHandle).Methods("GET")

	router.HandleFunc("/api/session/new", env.RegistrationHandle).Methods("POST")

	router.HandleFunc("/api/session", env.LoginHandle).Methods("POST")
	router.HandleFunc("/api/vkauth", env.VKRegister)
	return router
}
