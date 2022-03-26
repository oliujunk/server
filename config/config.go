package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type GlobalConfig struct {
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"database"`
	ApiServer struct {
		Port int `json:"port"`
	} `json:"apiServer"`
	CommandServer struct {
		Port int `json:"port"`
	} `json:"commandServer"`
}

var GlobalConfiguration GlobalConfig

func init() {
	config, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Printf(err.Error())
	}
	err = json.Unmarshal(config, &GlobalConfiguration)
	if err != nil {
		log.Printf(err.Error())
	}
}
