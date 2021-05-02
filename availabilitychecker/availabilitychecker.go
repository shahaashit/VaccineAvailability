package availabilitychecker

import (
	"VaccineAvailability/config"
	"VaccineAvailability/log"
	"VaccineAvailability/models"
	"VaccineAvailability/throttle"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func GetDataForMultiplePincodes(pinCodes []string) (finalData models.CentersList) {
	th := throttle.NewThrottle(5)
	for _, pinCode := range pinCodes {
		pinCode := pinCode
		th.Do()
		go func() {
			defer th.Done()
			resp, err := getDataForPincode(pinCode)
			if err != nil {
				log.Error("error while getting data from api: ", err)
				return
			}
			validCenters := resp.Centers.GetValidCenters(config.Config.Age)
			finalData = append(finalData, validCenters...)
		}()
	}
	th.Finish()
	return
}

func getDataForPincode(pincode string) (*models.HttpResponse, error) {
	finalResp := &models.HttpResponse{}

	q := url.Values{}
	q.Set("pincode", pincode)
	q.Set("date", time.Now().Format("02-01-2006"))
	u := url.URL{
		Scheme:   "https",
		Host:     "cdn-api.co-vin.in",
		Path:     "api/v2/appointment/sessions/public/calendarByPin",
		RawQuery: q.Encode(),
	}
	urlToCall := u.String()
	log.Debug(urlToCall)

	response, err := http.Get(urlToCall)
	if err != nil {
		return nil, errors.New("error while calling url:" + err.Error())
	}
	apResponseInBytes, err := ioutil.ReadAll(response.Body)
	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Error("error while closing body : ", err)
		}
	}()

	err = json.Unmarshal(apResponseInBytes, &finalResp)
	return finalResp, err
}
