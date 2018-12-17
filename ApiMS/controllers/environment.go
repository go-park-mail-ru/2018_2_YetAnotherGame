package controllers

import (
	"github.com/2018_2_YetAnotherGame/ApiMS/middlewares"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type Environment struct {
	DB      *gorm.DB
	Log     *middlewares.AccessLogger
	Counter *prometheus.CounterVec
	Conn    *grpc.ClientConn
}

func (env *Environment) InitDB(dialect, connStr string) {
	db, err := gorm.Open(dialect, connStr)
	if err != nil {
		logrus.Error(err)
	}
	env.DB = db
}

func (env *Environment) InitGrpc(port string) {
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		logrus.Fatalf("did not connect: %s", err)
	}

	env.Conn = conn
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
