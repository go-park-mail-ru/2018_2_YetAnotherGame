package main

import (
	"GameWS/game"
	"GameWS/middlewares"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

func Test(g *game.Game, w http.ResponseWriter, r *http.Request) {
	log.Printf("open connection")
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("cannot upgrade connection: ", err)
	}
	g.Register <- conn
}

func main() {
	addr := "localhost"
	port := 8000
	logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true})
	logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
		"host":   addr,
		"port":   port,
	}).Info("Starting server")
	AccessLogOut := new(middlewares.AccessLogger)
	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"http://127.0.0.1:3000"},                           // All origins
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"}, // Allowing only get, just an example
	})
	g := game.New()
	go g.Run()
	routerAuth := mux.NewRouter()
	routerAuth.HandleFunc("/ws", func(output http.ResponseWriter, request *http.Request) {
		Test(g, output, request)
	})
	authHandler := middlewares.AuthMiddleware(routerAuth)
	contextLogger := logrus.WithFields(logrus.Fields{
		"mode":   "[access_log]",
		"logger": "LOGRUS",
	})
	logrus.SetFormatter(&logrus.JSONFormatter{})
	AccessLogOut.LogrusLogger = contextLogger
	siteHandler := AccessLogOut.AccessLogMiddleware(authHandler)
	siteHandler = middlewares.PanicMiddleware(siteHandler)
	handler := c.Handler(siteHandler)
	http.ListenAndServe(":8081", handler)
}
