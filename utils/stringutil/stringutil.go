package stringutil

import (
	"regexp"
)

func CheckStringExistsInSlice(haystack []string, needle string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

func MatchRegex(pattern string, subject string) bool {
	match, err := regexp.Match(pattern, []byte(subject))
	if err != nil {
		return false
	}
	return match
}

func UniqueValues(intSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
