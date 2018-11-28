package middlewares

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func  PanicMiddleware(next http.Handler) http.Handler {
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
