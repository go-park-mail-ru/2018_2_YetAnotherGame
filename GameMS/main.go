package main

import (
	"GameMS/game"
	"GameMS/middlewares"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

func  Test(m *game.Metrics, g *game.Game, w http.ResponseWriter, r *http.Request) {
	log.Printf("open connection")
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("cannot upgrade connection: ", err)

	}
	g.Register <- conn
	m.Counter.WithLabelValues("rooms").Observe(float64(g.AmountRooms))
}

func main() {
	addr := "localhost"
	port := 8000
	met:=game.Metrics{}
	met.Counter=prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:"method_counter",
		Help:"count",
	},
		[]string{"rooms"},
	)
	prometheus.MustRegister(met.Counter)
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
		Test(&met, g, output, request)
	})
	conn, err := grpc.Dial(":7777", grpc.WithInsecure())
	if err != nil {
		logrus.Fatalf("did not connect: %s", err)
	}
	authHandler := middlewares.AuthMiddleware(routerAuth, conn)
	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	router.Handle("/ws", authHandler)
	contextLogger := logrus.WithFields(logrus.Fields{
		"mode":   "[access_log]",
		"logger": "LOGRUS",
	})



	logrus.SetFormatter(&logrus.JSONFormatter{})
	AccessLogOut.LogrusLogger = contextLogger
	siteHandler := AccessLogOut.AccessLogMiddleware(router)
	siteHandler = middlewares.PanicMiddleware(siteHandler, &met)
	handler := c.Handler(siteHandler)
	http.ListenAndServe(":8081", handler)
}
