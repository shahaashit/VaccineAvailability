package notifier

import (
	"VaccineAvailability/availabilitychecker/models"
	"VaccineAvailability/config"
	"VaccineAvailability/log"
	"encoding/json"
)

func Notify(finalData models.CentersList) {
	notificationList := make(map[string]models.CentersList)
	sortedData := finalData.GroupByPincode()
	for _, v := range config.AppConfiguration.NotificationList {
		var currentCentersList models.CentersList
		pincodes := v.GetPincodeSlice()
		for _, pincode := range pincodes {
			centersListForCurrentPincode := sortedData[pincode].GetAvailableCentersForAge(v.Age)
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
		err = sendMail(emailId, "vaccine availability notifier", string(byteData))
		if err != nil {
			log.Error("error while sending mail: ", err)
		}
	}
}
