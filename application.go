package main

import (
	"VaccineAvailability/availabilitychecker"
	"VaccineAvailability/availabilitychecker/models"
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
	"strings"
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

	_, err := c.AddFunc("@every "+config.AppConfiguration.CronFreq, process)
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
		districtId := r.URL.Query().Get("districtid")
		pincode := r.URL.Query().Get("pincode")
		var (
			resp *models.HttpResponse
			err  error
		)
		if districtId != "" {
			resp, err = availabilitychecker.GetDataForDistrictId(districtId)
		} else if pincode != "" {
			resp, err = availabilitychecker.GetDataForPincode(pincode)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("incorrect input values"))
			return
		}
		if err != nil {
			log.Error("request: error while getting data from service: ", err)
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

func process() {
	if killSwitch {
		return
	}
	hourId := utils.GetCurrentIstTime().Hour()
	if hourId >= 23 || hourId <= 7 { //throttle for non peak hours
		if utils.GetCurrentIstTime().Minute()%5 != 0 {
			log.Error("returning from processing in non-peak hours")
			return
		}
	}

	for _, channelConfig := range config.AppConfiguration.DistrictConfig.ChannelConfigList {
		var centersList models.CentersList
		var pincodeList []string

		if channelConfig.Pincode != "" {
			pincodeList = strings.Split(channelConfig.Pincode, ",")
		}
		if channelConfig.DistrictId != "" {
			resp, err := availabilitychecker.GetDataForDistrictId(channelConfig.DistrictId)
			if err != nil {
				log.Error("error while getting data from service: ", err)
				continue
			}
			centersList = resp.Centers

			if len(pincodeList) > 0 {
				centersList = centersList.FilterForPincodes(pincodeList)
			}
		} else if len(pincodeList) > 0 {
			centersList = availabilitychecker.GetDataForMultiplePincodes(pincodeList)
		} else {
			continue
		}

		centersList = centersList.GetFilteredCenters(true, channelConfig.AgeGroup)
		log.Debug("final selected centers: ", centersList)
		if len(centersList) > 0 {
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
