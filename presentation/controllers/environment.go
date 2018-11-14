package controllers

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type Environment struct {
	DB  *gorm.DB
	Log *logrus.Logger
}

func (env *Environment) InitDB(dialect, connStr string) {
	db, err := gorm.Open(dialect, connStr)
	if err != nil {
		logrus.Error(err)
	}
	env.DB = db
}

// func (env *Environment) initLog() {

// }
