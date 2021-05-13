package notifier

import (
	"VaccineAvailability/availabilitychecker/models"
	"VaccineAvailability/config"
	"VaccineAvailability/log"
	"VaccineAvailability/notifier/mail"
	"VaccineAvailability/notifier/telegram"
	"encoding/json"
	"strconv"
)

func NotifyUsingTelegram(channelConf config.TelegramChannelsConf, finalData models.CentersList) {
	if !currentSession.isCurrentDataNew(channelConf, finalData) { //throttle sending notification if data hasn't changed
		log.Info("not sending notification for district: ", channelConf.DistrictId)
		return
	}

	telegramData := parseDataForTelegramMessage(finalData)
	failures := false
	for _, data := range telegramData {
		err := telegram.SendTelegramMessage(channelConf.ChannelId, data)
		if err != nil {
			failures = true
			log.Error("error while sending telegram message: ", err)
		}
		log.Info("sending data to telegram channel for district: " + channelConf.DistrictId)
	}

	if !failures {
		currentSession.saveCurrentData(channelConf, finalData)
	}
}

func NotifyUsingMailForPincodeConfig(finalData models.CentersList) {
	notificationList := make(map[string]models.CentersList)
	sortedData := finalData.GroupByPincode()
	for _, v := range config.AppConfiguration.PincodeConfigList {
		var currentCentersList models.CentersList
		pincodes := v.GetPincodeSlice()
		for _, pincode := range pincodes {
			centersListForCurrentPincode := sortedData[pincode].GetFilteredCenters(true, v.AgeGroup)
			if _, ok := sortedData[pincode]; ok {
				currentCentersList = append(currentCentersList, centersListForCurrentPincode...)
			}
		}
		if len(currentCentersList) > 0 {
			if notificationList[v.Email] != nil {
				notificationList[v.Email] = append(notificationList[v.Email], currentCentersList...)
			} else {
				notificationList[v.Email] = currentCentersList
			}
		}
	}

	log.Debug(notificationList)

	for emailId, data := range notificationList {
		byteData, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			log.Error("error while marshalling data: ", err)
		}
		err = mail.SendMail(emailId, "vaccine availability notifier", string(byteData))
		if err != nil {
			log.Error("error while sending mail: ", err)
		}
		log.Info("sending mail to : ("+emailId+") for data: ", data)
	}
}

func parseDataForTelegramMessage(list models.CentersList) []string {
	var result []string
	tempResult := ""
	for i, center := range list {
		temp := strconv.Itoa(i+1) + ". <strong><u>" + center.Name + "</u></strong>\n" +
			"District: " + center.Districtname + "\n" +
			"Ward: " + center.BlockName + "\n" +
			"Pincode: " + strconv.Itoa(center.Pincode) + "\n" +
			"Fee Type: " + center.FeeType + "\n"
		for _, session := range center.Sessions {
			ageValue := "18-44"
			if session.MinAgeLimit == 45 {
				ageValue = "45+"
			}
			temp += session.Date + ": " + strconv.Itoa(session.AvailableCapacity) + " slots available for vaccine " + session.Vaccine + " for age group - " + ageValue + "\n"
		}
		temp += "\n\n"

		if len(tempResult+temp) > 4000 { //check for data length
			result = append(result, tempResult)
			tempResult = ""
		}
		tempResult += temp
	}
	if tempResult != "" {
		result = append(result, tempResult)
		tempResult = ""
	}
	return result
}
