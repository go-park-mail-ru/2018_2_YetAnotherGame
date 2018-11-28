package middlewares

import (
	"2018_2_YetAnotherGame/grpcModules"
	"2018_2_YetAnotherGame/infostructures/functions"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		logrus.Info("authMiddleware", r.URL.Path)
		session, err := r.Cookie("sessionid")
		if err != nil {
			logrus.Info("Unauthorized")
			err := functions.SendStatus("Unauthorized", w, http.StatusUnauthorized)
			if err != nil {
				logrus.Error(err)
			}
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
