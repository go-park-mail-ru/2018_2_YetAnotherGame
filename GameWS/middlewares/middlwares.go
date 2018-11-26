package middlewares

import (
	"2018_2_YetAnotherGame/grpcModules"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type AccessLogger struct {
	LogrusLogger *logrus.Entry
}

func (ac *AccessLogger) AccessLogMiddleware(next http.Handler) http.Handler {
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

func PanicMiddleware(next http.Handler) http.Handler {
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

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("auth")
		logrus.Info("authMiddleware", r.URL.Path)
		session, err := r.Cookie("sessionid")
		if err != nil {
			logrus.Info("Unauthorized")
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}
		id := session.Value

		status := grpcModules.SendCheckInfo(id)
		fmt.Println(status)
		if status == "Unauthorized" {
			logrus.Info("Unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}
