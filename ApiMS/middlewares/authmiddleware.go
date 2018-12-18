package middlewares

import (
	"github.com/go-park-mail-ru/2018_2_YetAnotherGame/ApiMS/grpcModules"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func AuthMiddleware(next http.Handler, conn *grpc.ClientConn) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("authMiddleware", r.URL.Path)
		session, err := r.Cookie("sessionid")
		if err != nil {
			logrus.Info("Unauthorized")
			//err := functions.SendStatus("Unauthorized", w, http.StatusUnauthorized)
			//if err != nil {
			//	logrus.Error(err)
			//}
			return
		}
		id := session.Value
		status := grpcModules.SendCheckInfo(id, conn)
		fmt.Println(status)
		if status == "Unauthorized" {
			logrus.Info("Unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}
