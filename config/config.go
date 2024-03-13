package config

import (
	"currentPrice/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	Host      string
	Port      string
	DBname    string
	User      string
	Pass      string
	BaseURL   string
	Leituras  int
	TimeSleep int64
	Config    *models.ConfigStruct
)

func ReadFile() {
	file, err := ioutil.ReadFile("config.json")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = json.Unmarshal(file, &Config)

	Host = Config.Host
	Port = Config.Port
	DBname = Config.DBname
	User = Config.User
	Pass = Config.Pass
	BaseURL = Config.BaseURL
	TimeSleep = Config.TimeSleep
	Leituras = Config.Leituras
}
