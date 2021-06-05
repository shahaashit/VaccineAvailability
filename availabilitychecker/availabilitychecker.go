package availabilitychecker

import (
	"VaccineAvailability/availabilitychecker/models"
	"VaccineAvailability/log"
	"VaccineAvailability/utils"
	"VaccineAvailability/utils/throttle"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GetDataForMultiplePincodes(pinCodes []string) (finalData models.CentersList) {
	th := throttle.NewThrottle(10)
	for _, pinCode := range pinCodes {
		pinCode := pinCode
		th.Do()
		go func() {
			defer th.Done()
			resp, err := GetDataForPincode(pinCode)
			if err != nil {
				log.Error("error while getting data from pincode api: ", err)
				return
			}
			finalData = append(finalData, resp.Centers...)
		}()
	}
	th.Finish()
	return
}

func GetDataForPincode(pincode string) (*models.HttpResponse, error) {
	finalResp := &models.HttpResponse{}

	q := url.Values{}
	q.Set("pincode", pincode)
	q.Set("date", utils.GetCurrentIstTime().Format("02-01-2006"))
	u := url.URL{
		Scheme:   "https",
		Host:     "cdn-api.co-vin.in",
		Path:     "api/v2/appointment/sessions/public/calendarByPin",
		RawQuery: q.Encode(),
	}
	urlToCall := u.String()
	log.Debug("current time used: ", utils.GetCurrentIstTime())
	log.Debug(urlToCall)

	var req *http.Request
	req, err := http.NewRequest("GET", urlToCall, nil)
	if err != nil {
		log.Error("error while making request : ", err)
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Accept-Language", "hi_IN")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36")

	var defaultClient = &http.Client{}
	response, err := defaultClient.Do(req)
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

func GetDataForDistrictId(districtId string) (*models.HttpResponse, error) {
	finalResp := &models.HttpResponse{}

	q := url.Values{}
	q.Set("district_id", districtId)
	q.Set("date", utils.GetCurrentIstTime().Format("02-01-2006"))
	u := url.URL{
		Scheme:   "https",
		Host:     "cdn-api.co-vin.in",
		Path:     "api/v2/appointment/sessions/public/calendarByDistrict",
		RawQuery: q.Encode(),
	}
	urlToCall := u.String()
	log.Debug("current time used: ", utils.GetCurrentIstTime())
	log.Debug(urlToCall)

	var req *http.Request
	req, err := http.NewRequest("GET", urlToCall, nil)
	if err != nil {
		log.Error("error while making request : ", err)
		return nil, err
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Accept-Language", "hi_IN")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36")

	var defaultClient = &http.Client{}
	response, err := defaultClient.Do(req)
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
