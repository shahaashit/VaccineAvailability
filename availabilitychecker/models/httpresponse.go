package models

import (
	"VaccineAvailability/utils"
	"sort"
	"strconv"
)

type HttpResponse struct {
	Centers CentersList `json:"centers"`
}

type Centers struct {
	CenterId     int           `json:"center_id"`
	Name         string        `json:"name"`
	StateName    string        `json:"state_name"`
	Districtname string        `json:"district_name"`
	BlockName    string        `json:"block_name"`
	Pincode      int64         `json:"pincode"`
	Lat          int64         `json:"lat"`
	Long         int64         `json:"long"`
	From         string        `json:"from"`
	To           string        `json:"to"`
	FeeType      string        `json:"fee_type"`
	Sessions     []Sessions    `json:"sessions"`
	VaccineFees  []VaccineFees `json:"vaccine_fees"`
}

type Sessions struct {
	SessionId         string   `json:"session_id"`
	Date              string   `json:"date"`
	AvailableCapacity int      `json:"available_capacity"`
	MinAgeLimit       int      `json:"min_age_limit"`
	Vaccine           string   `json:"vaccine"`
	Slots             []string `json:"slots"`
}

type VaccineFees struct {
	Vaccine string `json:"vaccine"`
	Fee     string `json:"fee"`
}

type CentersList []Centers

func (c CentersList) GetFilteredCenters(onlyAvailableCheck bool, currentAge int) CentersList {
	var result CentersList
	for _, center := range c {
		var validSessions []Sessions
		for _, session := range center.Sessions {
			if onlyAvailableCheck && session.AvailableCapacity <= 0 {
				continue
			}
			if currentAge != -1 && session.MinAgeLimit != currentAge {
				continue
			}
			validSessions = append(validSessions, session)
		}
		if len(validSessions) > 0 {
			center.Sessions = validSessions
			result = append(result, center)
		}
	}
	return result
}

func (c CentersList) GetAllSortedCenterIds() []int {
	var result []int
	for _, center := range c {
		result = append(result, center.CenterId)
	}
	result = utils.UniqueValues(result)
	sort.Ints(result)
	return result
}

func (c CentersList) GroupByPincode() map[string]CentersList {
	finalData := make(map[string]CentersList)
	for _, center := range c {
		pincode := strconv.FormatInt(center.Pincode, 10)
		finalData[pincode] = append(finalData[pincode], center)
	}
	return finalData
}
