package routes

import (
	"net/http"

	"2018_2_YetAnotherGame/presentation/controllers"

	"github.com/gorilla/mux"
)

func Router(env *controllers.Environment) http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/api/vkauth", env.VKRegister)
	router.HandleFunc("/api/session", env.LoginHandle).Methods("POST")
	router.HandleFunc("/api/session/new", env.RegistrationHandle).Methods("POST")
	router.HandleFunc("/api/users/me", env.MeHandle).Methods("GET")
	router.HandleFunc("/api/leaders", env.ScoreboardHandle).Methods("GET")
	router.HandleFunc("/api/session", env.LogOutHandle).Methods("DELETE")
	router.HandleFunc("/api/avatar", env.AvatarHandle)
	return router
}
