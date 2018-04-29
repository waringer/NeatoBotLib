package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"../NeatoBotLib"
)

type configuration struct {
	MetaURL  string `json:"URL"`
	EMail    string `json:"eMail"`
	Password string `json:"password"`
}

var Conf configuration

func main() {
	confFile := flag.String("c", "sample.conf", "config file to use")
	flag.Parse()
	
	err := loadConfig(*confFile)
	if err != nil {
		log.Println("can't read conf file", *confFile)
		saveConfig(*confFile)
	}

	auth := NeatoBotLib.Auth(Conf.MetaURL, Conf.EMail, Conf.Password)
	dash := NeatoBotLib.GetDashboard(Conf.MetaURL, auth)

	for _, rob := range dash.Robots {
		state := NeatoBotLib.GetRobotState(auth, rob)
		fmt.Println(state)
	}
}

func loadConfig(filename string) error {
	DefaultConf := configuration{MetaURL: "", EMail: "", Password: ""}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		Conf = configuration{}
		return err
	}

	err = json.Unmarshal(bytes, &DefaultConf)
	if err != nil {
		Conf = configuration{}
		return err
	}

	Conf = DefaultConf
	return nil
}

func saveConfig(filename string) error {
	bytes, err := json.MarshalIndent(Conf, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, bytes, 0644)
}
