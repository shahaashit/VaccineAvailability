package config

import (
	"VaccineAvailability/utils/stringutil"
	"strings"
)

var AppConfiguration Conf

type Conf struct {
	MailClientConfig  MailClientConf `yaml:"mailclientconf"`
	PincodeConfigList []PincodeConf  `yaml:"pincodeconfig"`
	DistrictConfig    DistrictConf   `yaml:"districtconfig"`
	LogExtraInfo      bool           `yaml:"logextrainfo"`
	CronFreq          string         `yaml:"cronfreq"`
	TestEnv           bool           `yaml:"testenv"`
}

type DistrictConf struct {
	Token                 string                 `yaml:"token"`
	TestChannelConfigList []TelegramChannelsConf `yaml:"testchannelconfig"`
	ChannelConfigList     []TelegramChannelsConf `yaml:"channelconfig"`
}

type TelegramChannelsConf struct {
	DistrictId string `yaml:"districtid"`
	Pincode    string `yaml:"pincode"`
	ChannelId  string `yaml:"channelid"`
	AgeGroup   int    `yaml:"agegroup"`
}

type MailClientConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
}

type PincodeConf struct {
	Email    string `yaml:"email"`
	Pincodes string `yaml:"pincode"`
	AgeGroup int    `yaml:"agegroup"`
}

func (n PincodeConf) GetPincodeSlice() []string {
	return strings.Split(n.Pincodes, ",")
}

func (c Conf) GetUniquePincodesFromAllConfigs() []string {
	var returnVal []string
	for _, conf := range c.PincodeConfigList {
		returnVal = append(returnVal, conf.GetPincodeSlice()...)
	}
	returnVal = stringutil.UniqueValues(returnVal)
	return returnVal
}
