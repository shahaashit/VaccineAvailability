package main

import (
	"VaccineAvailability/availabilitychecker"
	"VaccineAvailability/config"
	_ "VaccineAvailability/config"
	"VaccineAvailability/log"
	"VaccineAvailability/notifier"
	"github.com/robfig/cron/v3"
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
	_, err := c.AddFunc("@every "+config.AppConfiguration.CronFreq, process)
	if err != nil {
		log.Fatal("error while setting up cron: ", err)
	}

	c.Start()
}

func process() {
	resp := availabilitychecker.GetDataForMultiplePincodes(config.AppConfiguration.GetUniquePincodes())
	if len(resp) > 0 {
		notifier.Notify(resp)
	}
}
