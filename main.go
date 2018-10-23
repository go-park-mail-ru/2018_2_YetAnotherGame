package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"
	"log"
	"os"
	"time"

	"2018_2_YetAnotherGame/handlers"
	"2018_2_YetAnotherGame/models"

	"github.com/gorilla/mux"
	//"log"
	"net/http"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type AccessLogger struct {
	StdLogger    *log.Logger

}
func (ac *AccessLogger) accessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		fmt.Printf("FMT [%s] %s, %s %s\n",
			r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))

		log.Printf("LOG [%s] %s, %s %s\n",
			r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))

		ac.StdLogger.Printf("[%s] %s, %s %s\n",
			r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))


		})
	}



func main() {
	AccessLogOut := new(AccessLogger)

	// std
	AccessLogOut.StdLogger = log.New(os.Stdout, "STD ", log.LUTC|log.Lshortfile)

	db, err := gorm.Open("postgres", "host=127.0.0.1 port=5432 user=test_user dbname=backend password=1")
	fmt.Println(err)
	defer db.Close()
	//env:=&database.Env{db}
	ids := make(map[string]string, 0)
	//users := make(map[string]*models.User, 0)
	users:=models.UsersMap{}
	users.Const()

	//user:=new(user){"a@a","f1","l1","u1","qwerty",5}
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Session{})
	q1 := models.User{"1","af@a", "f1", "l1", "u1", "qwerty", 5, ""}
	q2 := models.User{"2","asf@a", "f1", "l1", "u1", "qwerty", 6, ""}
	q3 := models.User{"3","asfg@a", "f1", "l1", "u1", "qwerty", 54, ""}
	q4 := models.User{"4","asdg@a", "f1", "l1", "u1", "qwerty", 7, ""}
	q5 := models.User{"5","asdg@a", "f1", "l1", "u1", "qwerty", 6, ""}
	q6 := models.User{"6","asdg@a", "f1", "l1", "u1", "qwerty", 9, ""}
	db.Create(&q1)
	db.Create(&q2)
	db.Create(&q3)
	db.Create(&q4)
	db.Create(&q5)
	db.Create(&q6)


	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"http://127.0.0.1:3000"},                           // All origins
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"}, // Allowing only get, just an example

	})


	//mux := http.NewServeMux()
	router := mux.NewRouter()

	router.HandleFunc("/api/user", func(output http.ResponseWriter, request *http.Request) {
		handlers.Leaders(db, output, request)
	}).Methods("GET")

	router.HandleFunc("/api/session/new", func(output http.ResponseWriter, request *http.Request) {
		handlers.SignUp(db,ids, users, output, request)
	}).Methods("POST")

	router.HandleFunc("/api/session", func(output http.ResponseWriter, request *http.Request) {
		handlers.Login(db, output, request)
	}).Methods("POST")

	router.HandleFunc("/api/user/me", func(output http.ResponseWriter, request *http.Request) {
		handlers.Me(db, output, request)
	}).Methods("GET")

	router.HandleFunc("/api/session", handlers.Logout).Methods("DELETE")

	router.HandleFunc("/api/user/me", func(output http.ResponseWriter, request *http.Request) {
		handlers.Update(db, output, request)
	}).Methods("POST")

	router.HandleFunc("/api/upload", func(output http.ResponseWriter, request *http.Request) {
		handlers.Upload(db, output, request)
	}).Methods("POST")

	http.Handle("/", router)

	fmt.Println("Server listening port 8000")
	//log.Fatal(http.ListenAndServe(":8000", c.Handler(router)))
	siteHandler := AccessLogOut.accessLogMiddleware(router)
	handler := c.Handler(siteHandler)
	http.ListenAndServe(":8000", handler)
}
