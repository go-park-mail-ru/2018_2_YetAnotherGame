package main

import (
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/FogCreek/mini"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"

	"2018_2_YetAnotherGame/handlers"
	"2018_2_YetAnotherGame/models"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	//"log"
	"net/http"
)

type AccessLogger struct {
	LogrusLogger *logrus.Entry
}

func (ac *AccessLogger) accessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		ac.LogrusLogger.WithFields(logrus.Fields{
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
			"work_time":   time.Since(start),
		}).Info(r.URL.Path)
	})
}

func panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("panicMiddleware", r.URL.Path)
		defer func() {
			if err := recover(); err != nil {
				logrus.Panic("recovered", err)
				http.Error(w, "Internal server error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(db *gorm.DB,next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("authMiddleware", r.URL.Path)
		id, err := r.Cookie("sessionid")
		id2 := id.Value
		if err != nil {
			logrus.Error("no auth at", r.URL.Path)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		tmp := models.Session{}
		db.Table("sessions").Select("id, email").Where("id = ?", id2).Scan(&tmp)
		if tmp.Email== "" {
			logrus.Error("no auth at", r.URL.Path)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func params() string {
	u, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}
	pwd, _ := os.Getwd()
	cfg, err := mini.LoadConfiguration(pwd + "/config/DBsettings.txt")
	if err != nil {
		logrus.Error(err)
	}

	info := fmt.Sprintf("host=%s port=%s dbname=%s "+
		"sslmode=%s user=%s password=%s ",
		cfg.String("host", "127.0.0.1"),
		cfg.String("port", "5432"),
		cfg.String("dbname", u.Username),
		cfg.String("sslmode", "disable"),
		cfg.String("user", u.Username),
		cfg.String("pass", ""),
	)
	return info
}

func main() {
	addr := "localhost"
	port := 8000
	// logrus
	logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true})
	logrus.WithFields(logrus.Fields{

		"logger": "LOGRUS",
		"host":   addr,
		"port":   port,

	}).Info("Starting server")
	AccessLogOut := new(AccessLogger)

	db, err := gorm.Open("postgres", params())
	if err != nil {
		logrus.Error(err)
	}
	defer db.Close()

	//test users to fill the db
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Session{})
	//q1 := models.User{"1","af@a", "f1", "l1", "u1", "qwerty", 5, ""}
	//q2 := models.User{"2","asf@a", "f1", "l1", "u1", "qwerty", 6, ""}
	//q3 := models.User{"3","asfg@a", "f1", "l1", "u1", "qwerty", 54, ""}
	//q4 := models.User{"4","asdg@a", "f1", "l1", "u1", "qwerty", 7, ""}
	//q5 := models.User{"5","asdg@a", "f1", "l1", "u1", "qwerty", 6, ""}
	//q6 := models.User{"6","asdg@a", "f1", "l1", "u1", "qwerty", 9, ""}
	//db.Create(&q1)
	//db.Create(&q2)
	//db.Create(&q3)
	//db.Create(&q4)
	//db.Create(&q5)
	//db.Create(&q6)

	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"http://127.0.0.1:3000"},                           // All origins
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"}, // Allowing only get, just an example

	})

	//mux := http.NewServeMux()
	routerAuth := mux.NewRouter()
	routerAuth.HandleFunc("/api/user/me", func(output http.ResponseWriter, request *http.Request) {
		handlers.Me(db, output, request)
	}).Methods("GET")
	routerAuth.HandleFunc("/api/user/me", func(output http.ResponseWriter, request *http.Request) {
		handlers.Update(db, output, request)
	}).Methods("POST")
	routerAuth.HandleFunc("/api/upload", func(output http.ResponseWriter, request *http.Request) {
		handlers.Upload(db, output, request)
	}).Methods("POST")
	routerAuth.HandleFunc("/api/session", handlers.Logout).Methods("DELETE")
	authHandler := authMiddleware(db, routerAuth)

	router := mux.NewRouter()
	router.Handle("/api/user/me",authHandler)
	router.Handle("/api/session",authHandler).Methods("DELETE")
	router.Handle("/api/upload",authHandler)
	router.HandleFunc("/api/leaders", func(output http.ResponseWriter, request *http.Request) {
		handlers.Leaders(db, output, request)
	}).Methods("GET")

	router.HandleFunc("/api/session/new", func(output http.ResponseWriter, request *http.Request) {
		handlers.SignUp(db, output, request)
	}).Methods("POST")

	router.HandleFunc("/api/session", func(output http.ResponseWriter, request *http.Request) {
		handlers.Login(db, output, request)
	}).Methods("POST")


	http.Handle("/", router)
	// logrus
	contextLogger := logrus.WithFields(logrus.Fields{
		"mode":   "[access_log]",
		"logger": "LOGRUS",
	})
	logrus.SetFormatter(&logrus.JSONFormatter{})
	AccessLogOut.LogrusLogger = contextLogger
	//fmt.Println("Server listening port 8000")
	//log.Fatal(http.ListenAndServe(":8000", c.Handler(router)))

	siteHandler := AccessLogOut.accessLogMiddleware(router)
	siteHandler = panicMiddleware(siteHandler)
	handler := c.Handler(siteHandler)
	http.ListenAndServe(":8000", handler)
}
