package config

import (
	"VaccineAvailability/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

var Config Conf

type Conf struct {
	Pincode string `yaml:"pincode"`
	Age     int64  `yaml:"age"`
	Debug   bool   `yaml:"debug"`
}

func init() {
	SetConfig()
}

func SetConfig() {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("config read file error: " + err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &Config)
	if Config.Pincode == "" {
		log.Fatal("Incorrect pincode value in config")
	}
	if Config.Age < 0 {
		log.Fatal("Incorrect age value in config")
	}
	log.SetWriter(os.Stdout, Config.Debug)
}
