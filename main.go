package main

import (
	"VaccineAvailability/availabilitychecker"
	"VaccineAvailability/config"
	"VaccineAvailability/log"
	"github.com/robfig/cron/v3"
	"strings"
)

var (
	shutDownCh = make(chan bool)
)

func main() {
	go startCron()
	<-shutDownCh //waiting for this to never stop the service
}

func startCron() {
	c := cron.New()
	//c.AddFunc("@every 5s", config.SetConfig)
	_, err := c.AddFunc("@every 1m", process)
	if err != nil {
		log.Fatal("error while setting up cron: ", err)
	}

	c.Start()
}

func process() {
	resp := availabilitychecker.GetDataForMultiplePincodes(strings.Split(config.Config.Pincode, ","))
	if len(resp) > 0 {
		log.Info(resp.GroupByPincode())
	}
}
