package utils

import "time"

var (
	istLoc *time.Location
)

func init() {
	istLoc, _ = time.LoadLocation("Asia/Kolkata")
}

func GetCurrentIstTime() time.Time {
	return time.Now().In(istLoc)
}

func UniqueValues(intSlice []int) []int {
	keys := make(map[int]bool)
	var list []int
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
