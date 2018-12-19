package main

import (
	"github.com/go-park-mail-ru/2018_2_YetAnotherGame/ApiMS/controllers"
	"github.com/go-park-mail-ru/2018_2_YetAnotherGame/ApiMS/middlewares"
	"github.com/go-park-mail-ru/2018_2_YetAnotherGame/ApiMS/routes"
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
	_, err:=toml.DecodeFile("/home/ubuntu/go/src/config/DBsettings.toml", conf)
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Printf("%s", conf.String())
	return conf.String()
}

func main() {
	env := controllers.Environment{}
	env.InitLog()
	env.InitDB("postgres", dbSettings())
	env.InitGrpc(":8082")
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
	//http.ListenAndServe(":8000", r)
	err:=http.ListenAndServeTLS(":8081", "/etc/letsencrypt/live/yet-another-game.ml/fullchain.pem","/etc/letsencrypt/live/yet-another-game.ml/privkey.pem",r)
	if err!=nil{
		fmt.Println(err)
	}
}
