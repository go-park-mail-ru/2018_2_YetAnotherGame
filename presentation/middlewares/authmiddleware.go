package middlewares

import (
	"2018_2_YetAnotherGame/domain/models"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func AuthMiddleware(next http.Handler, db *gorm.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("authMiddleware", r.URL.Path)
		session, err := r.Cookie("sessionid")
		if err != nil {
			logrus.Info("Unauthorized")
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}
		id := session.Value
		if err != nil {
			logrus.Error("no auth at", r.URL.Path)
			http.Redirect(w, r, "/", http.StatusForbidden)
			return
		}

		tmp := models.Session{}
		db.Table("sessions").Select("id, email").Where("id = ?", id).Scan(&tmp)
		if tmp.Email == "" {
			logrus.Error("no auth at", r.URL.Path)
			http.Redirect(w, r, "/", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
