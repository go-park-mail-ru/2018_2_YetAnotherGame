package controllers

import (

	"2018_2_YetAnotherGame/presentation/middlewares"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type Environment struct {
	DB  *gorm.DB
	Log *middlewares.AccessLogger

}

func (env *Environment) InitDB(dialect, connStr string) {
	db, err := gorm.Open(dialect, connStr)
	if err != nil {
		logrus.Error(err)
	}
	env.DB = db
}



func (env *Environment) InitLog() {
	logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true})
	env.Log = new(middlewares.AccessLogger)
	logContext := logrus.WithFields(logrus.Fields{
		"mode":   "[access_log]",
		"logger": "LOGRUS",
	})

	env.Log.LogrusLogger = logContext
}
