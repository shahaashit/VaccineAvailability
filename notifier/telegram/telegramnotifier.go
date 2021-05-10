package telegram

import (
	"VaccineAvailability/config"
	"VaccineAvailability/log"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type TelegramResponse struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

func SendTelegramMessage(channelId string, data string) error {
	q := url.Values{}
	q.Set("parse_mode", "HTML")
	q.Set("chat_id", channelId)
	q.Set("text", data)
	u := url.URL{
		Scheme:   "https",
		Host:     "api.telegram.org",
		Path:     "/bot" + config.AppConfiguration.DistrictConfig.Token + "/sendMessage",
		RawQuery: q.Encode(),
	}
	urlToCall := u.String()

	resp, err := http.PostForm(urlToCall, url.Values{})
	if err != nil {
		log.Error("error while sending data via telegram: ", err)
		return err
	}

	res := TelegramResponse{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	log.Info(res)
	if err != nil {
		//log.Error("error while decoding telegram message response: ", err)
		return err
	}

	if !res.Ok {
		//log.Error("error response from telegram while sending message: ", err)
		return errors.New(res.Description)
	}
	return nil
}
