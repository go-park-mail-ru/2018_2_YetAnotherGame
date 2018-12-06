package middlewares

import (
	"2018_2_YetAnotherGame/grpcModules"
	"google.golang.org/grpc"

	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

func AuthMiddleware(next http.Handler, conn *grpc.ClientConn) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		logrus.Info("authMiddleware", r.URL.Path)
		session, err := r.Cookie("sessionid")
		if err != nil {
			logrus.Info("Unauthorized")
			return
		}
		id := session.Value
		status:=grpcModules.SendCheckInfo(id, conn)
		fmt.Println(status)
		if status=="Unauthorized"{
			logrus.Info("Unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}
