package main

import (
	"fmt"

	"./handlers"
	"./models"

	"github.com/gorilla/mux"
	//"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	ids := make(map[string]string, 0)
	users := make(map[string]*models.User, 0)
	avatars := make(map[string]string, 0)
	//user:=new(user){"a@a","f1","l1","u1","qwerty",5}

	q1 := models.User{"af@a", "f1", "l1", "u1", "qwerty", 5}
	q2 := models.User{"asf@a", "f1", "l1", "u1", "qwerty", 6}
	q3 := models.User{"asfg@a", "f1", "l1", "u1", "qwerty", 54}
	q4 := models.User{"asdg@a", "f1", "l1", "u1", "qwerty", 7}
	q5 := models.User{"asdg@a", "f1", "l1", "u1", "qwerty", 6}
	q6 := models.User{"asdg@a", "f1", "l1", "u1", "qwerty", 9}
	users["1"] = &q1
	users["2"] = &q2
	users["3"] = &q3
	users["4"] = &q4
	users["5"] = &q5
	users["6"] = &q6

	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"http://127.0.0.1:3000"},                           // All origins
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"}, // Allowing only get, just an example
	})

	//mux := http.NewServeMux()
	router := mux.NewRouter()

	router.HandleFunc("/user", func(output http.ResponseWriter, request *http.Request) {
		handlers.Leaders(users, output, request)
	}).Methods("GET")

	router.HandleFunc("/session/new", func(output http.ResponseWriter, request *http.Request) {
		handlers.SignUp(ids, users, output, request)
	}).Methods("POST")

	router.HandleFunc("/session", func(output http.ResponseWriter, request *http.Request) {
		handlers.Login(ids, users, output, request)
	}).Methods("POST")

	router.HandleFunc("/user/me", func(output http.ResponseWriter, request *http.Request) {
		handlers.Me(users, avatars, output, request)
	}).Methods("GET")

	router.HandleFunc("/session", handlers.Logout).Methods("DELETE")

	router.HandleFunc("/user/me", func(output http.ResponseWriter, request *http.Request) {
		handlers.Update(users, output, request)
	}).Methods("POST")

	// TODO: изменить по ресту
	router.HandleFunc("/upload", func(output http.ResponseWriter, request *http.Request) {
		handlers.Upload(avatars, output, request)
	})

	http.Handle("/", router)

	fmt.Println("Server is listening...")
	//log.Fatal(http.ListenAndServe(":8000", c.Handler(router)))
	handler := c.Handler(router)
	http.ListenAndServe(":8000", handler)
}
