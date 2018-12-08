package main

import (
	"2018_2_YetAnotherGame/ApiMS/controllers"
	"2018_2_YetAnotherGame/ApiMS/middlewares"
	"2018_2_YetAnotherGame/ApiMS/routes"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/BurntSushi/toml"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DbConfig struct {
	Host    string `toml:"host"`
	Port    string `toml:"port"`
	Sslmode string `toml:"sslmode"`
	Dbname  string `toml:"dbname"`
	User    string `toml:"user"`
	Pass    string `toml:"pass"`
}

func (db DbConfig) String() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s "+
		"sslmode=%s user=%s password=%s ",
		db.Host, db.Port, db.Dbname, db.Sslmode, db.User, db.Pass,
	)
}

func dbSettings() string {
	conf := &DbConfig{}
	toml.DecodeFile("./config/DBsettings.toml", conf)
	fmt.Printf("%s", conf.String())
	return conf.String()
}

func main() {
	env := controllers.Environment{}
	env.InitLog()
	env.InitDB("postgres", dbSettings())
	env.InitGrpc(":7777")
	env.Counter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "method_counter",
		Help: "counter",
	},
		[]string{"method", "status"},
	)
	r := routes.Router(&env)
	r = env.Log.AccessLogMiddleware(
		middlewares.PanicMiddleware(
			middlewares.CORSMiddleware(
				r,
			),
		),
	)
	http.ListenAndServe(":8000", r)
}
