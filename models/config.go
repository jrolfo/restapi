package models

import (
	"fmt"
	"time"

	"github.com/tkanos/gonfig"
)

//Configuration struct
type Configuration struct {
	JwtKey   string
	Expires  time.Duration
	Database string
	User     string
	Password string
	Server   string
	Port     string
}

//Config is exported
var Config = Configuration{}

//InitConfig function
func InitConfig() {
	err := gonfig.GetConf("./config.json", &Config)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
