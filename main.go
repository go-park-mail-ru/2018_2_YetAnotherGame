package main

import (
	"fmt"
	"net/http"

	"2018_2_YetAnotherGame/presentation/controllers"
	"2018_2_YetAnotherGame/presentation/routes"

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
	// env.initLog()
	env.InitDB("postgres", dbSettings())

	router := routes.Router(&env)
	http.ListenAndServe(":8000", router)
}
