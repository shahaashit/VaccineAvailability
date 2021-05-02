package models

import "strconv"

type HttpResponse struct {
	Centers CentersList `json:"centers"`
}

type Centers struct {
	CenterId     int64         `json:"center_id"`
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
	AvailableCapacity int64    `json:"available_capacity"`
	MinAgeLimit       int64    `json:"min_age_limit"`
	Vaccine           string   `json:"vaccine"`
	Slots             []string `json:"slots"`
}

type VaccineFees struct {
	Vaccine string `json:"vaccine"`
	Fee     string `json:"fee"`
}

type CentersList []Centers

func (c CentersList) GetValidCenters(currentAge int64) CentersList {
	var result CentersList
	for _, center := range c {
		var validSessions []Sessions
		for _, session := range center.Sessions {
			if session.AvailableCapacity > 0 && session.MinAgeLimit <= currentAge {
				validSessions = append(validSessions, session)
			}
		}
		if len(validSessions) > 0 {
			center.Sessions = validSessions
			result = append(result, center)
		}
	}
	return result
}

func (c CentersList) GroupByPincode() map[string][]Centers {
	finalData := make(map[string][]Centers)
	for _, center := range c {
		pincode := strconv.FormatInt(center.Pincode, 10)
		finalData[pincode] = append(finalData[pincode], center)
	}
	return finalData
}
