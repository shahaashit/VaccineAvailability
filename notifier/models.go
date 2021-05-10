package notifier

import (
	"VaccineAvailability/availabilitychecker/models"
	"VaccineAvailability/config"
	"VaccineAvailability/utils"
	"reflect"
	"strconv"
	"sync"
	"time"
)

var (
	currentSession            MessageFloodData
	repeatMessageIntervalMins = 60
)

type MessageFloodData struct {
	sync.RWMutex
	data map[string]MessageFlood
}

type MessageFlood struct {
	LastUpdatedTime time.Time
	NoOfHits        int
	Data            []int
}

func (mf MessageFloodData) getKeyForLocalSession(channelConf config.TelegramChannelsConf) string {
	return channelConf.DistrictId + "~~" + strconv.Itoa(channelConf.AgeGroup)
}

func (mf MessageFloodData) isCurrentDataNew(channelConf config.TelegramChannelsConf, centersList models.CentersList) bool {
	mf.RLock()
	defer mf.RUnlock()

	keyToUse := mf.getKeyForLocalSession(channelConf)
	if savedData, ok := mf.data[keyToUse]; ok {
		if reflect.DeepEqual(savedData.Data, centersList.GetAllSortedCenterIds()) && utils.GetCurrentIstTime().Before(savedData.LastUpdatedTime.Add(time.Duration(repeatMessageIntervalMins)*time.Minute)) {
			return false
		}
	}
	return true
}

func (mf *MessageFloodData) saveCurrentData(channelConf config.TelegramChannelsConf, centersList models.CentersList) {
	mf.Lock()
	defer mf.Unlock()

	keyToUse := mf.getKeyForLocalSession(channelConf)
	mf.data[keyToUse] = MessageFlood{
		LastUpdatedTime: utils.GetCurrentIstTime(),
		Data:            centersList.GetAllSortedCenterIds(),
	}
}

func init() {
	currentSession = MessageFloodData{}
	currentSession.data = make(map[string]MessageFlood)
}
