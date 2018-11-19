package middlewares

import (
	"2018_2_YetAnotherGame/grpcModules"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/jinzhu/gorm"
)

func AuthMiddleware(next http.Handler, db *gorm.DB) http.Handler {
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

		status:=grpcModules.SendCheckInfo(id)
		fmt.Println(status)
		if status=="Unauthorized"{
			logrus.Info("Unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}
