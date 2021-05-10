package main

import (
	"VaccineAvailability/availabilitychecker"
	"VaccineAvailability/config"
	_ "VaccineAvailability/config"
	"VaccineAvailability/log"
	"VaccineAvailability/notifier"
	"VaccineAvailability/utils"
	"encoding/json"
	"github.com/robfig/cron/v3"
	"net/http"
	"os"
	"strconv"
)

var (
	shutDownCh = make(chan bool)
	killSwitch = false
)

func main() {
	go startCron()
	startWebServer()
	<-shutDownCh
}

func startCron() {
	c := cron.New()

	_, err := c.AddFunc("@every "+config.AppConfiguration.CronFreq, processForDistrict)
	if err != nil {
		log.Fatal("error while setting up cron: ", err)
	}

	c.Start()
}

func startWebServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Infof("Serving data to %s...\n", r.RemoteAddr)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("successful"))
	})

	http.HandleFunc("/getData", func(w http.ResponseWriter, r *http.Request) {
		resp, err := availabilitychecker.GetDataForDistrictId(r.URL.Query().Get("districtId"))
		if err != nil {
			log.Error("request: error while getting data from api: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		} else {
			finalOutput := resp.Centers
			onlyAvailable := false
			age := -1
			if r.URL.Query().Get("onlyAvailable") == "1" {
				onlyAvailable = true
			}
			if r.URL.Query().Get("age") != "" {
				val, e := strconv.Atoi(r.URL.Query().Get("age"))
				if e == nil {
					age = val
				}
			}
			finalOutput = finalOutput.GetFilteredCenters(onlyAvailable, age)
			w.WriteHeader(http.StatusOK)
			jsonOutput, _ := json.Marshal(&finalOutput)
			w.Write(jsonOutput)
		}
	})

	http.HandleFunc("/killSwitch", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("value") == "1" {
			log.Infof("kill switch: enabled")
			killSwitch = true
		} else {
			log.Infof("kill switch: disabled")
			killSwitch = false
		}

		w.WriteHeader(http.StatusOK)
		response := "cron service running"
		if killSwitch {
			response = "cron service stopped"
		}
		w.Write([]byte(response))
	})

	log.Infof("Listening on port %s\n\n", port)
	http.ListenAndServe(":"+port, nil)
}

func processForDistrict() {
	if killSwitch {
		return
	}
	hourId := utils.GetCurrentIstTime().Hour()
	if hourId >= 21 || hourId <= 8 { //throttle for non peak hours
		if utils.GetCurrentIstTime().Minute()%5 != 0 {
			log.Error("returning from processing in non-peak hours")
			return
		}
	}

	for _, channelConfig := range config.AppConfiguration.DistrictConfig.ChannelConfigList {
		resp, err := availabilitychecker.GetDataForDistrictId(channelConfig.DistrictId)
		if err != nil {
			log.Error("error while getting data from api: ", err)
			return
		}
		centersList := resp.Centers
		centersList = centersList.GetFilteredCenters(true, channelConfig.AgeGroup)
		if len(resp.Centers) > 0 {
			notifier.NotifyUsingTelegram(channelConfig, centersList)
		}
	}
}

func processForPincode() {
	resp := availabilitychecker.GetDataForMultiplePincodes(config.AppConfiguration.GetUniquePincodesFromAllConfigs())
	if len(resp) > 0 {
		notifier.NotifyUsingMailForPincodeConfig(resp)
	}
}
