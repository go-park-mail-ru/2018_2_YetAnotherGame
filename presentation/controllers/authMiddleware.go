package controllers

import (
	"2018_2_YetAnotherGame/domain/models"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (env *Environment) AuthMiddleware(next http.Handler) http.Handler {
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
		env.DB.Table("sessions").Select("id, email").Where("id = ?", id2).Scan(&tmp)
		if tmp.Email == "" {
			logrus.Error("no auth at", r.URL.Path)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
