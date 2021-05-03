package config

import (
	"VaccineAvailability/log"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

func init() {
	SetConfig()
}

func SetConfig() {
	yamlFile, err := ioutil.ReadFile("resources/config.yaml") //todo: please configure this file according to sample_config.yaml
	if err != nil {
		log.Fatal("config read file error: " + err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &AppConfiguration)
	if err != nil {
		log.Fatal("Incorrect config: ", err)
	}
	err = validateConfig()
	if err != nil {
		log.Fatal("Incorrect config: ", err)
	}
	log.SetWriter(os.Stdout, AppConfiguration.Debug)
}

func validateConfig() error {
	return nil //todo: check if this needs to be done
}
