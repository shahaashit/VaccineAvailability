package config

import (
	"VaccineAvailability/utils/stringutil"
	"strings"
)

var AppConfiguration Conf

type Conf struct {
	EmailConfig      EmailConf          `yaml:"emailconfig"`
	NotificationList []NotificationList `yaml:"notificationconfig"`
	Debug            bool               `yaml:"debug"`
	CronFreq         string             `yaml:"cronfreq"`
}

type EmailConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
}

type NotificationList struct {
	Email    string `yaml:"email"`
	Pincodes string `yaml:"pincode"`
	Age      int64  `yaml:"minage"`
	Debug    bool   `yaml:"debug"`
}

func (n NotificationList) GetPincodeSlice() []string {
	return strings.Split(n.Pincodes, ",")
}

func (c Conf) GetUniquePincodes() []string {
	var returnVal []string
	for _, conf := range c.NotificationList {
		returnVal = append(returnVal, conf.GetPincodeSlice()...)
	}
	returnVal = stringutil.UniqueValues(returnVal)
	return returnVal
}
