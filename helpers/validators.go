package helpers

import (
	"time"
)

func ValidateDateQueryParam(queryParamDate string, dateFormat string) bool {

	_, err := time.Parse(dateFormat, queryParamDate)
	if err != nil {
		return false
	}

	return true
}

func ValidateDateRange(startDate string, endDate string, dateFormat string) bool {

	date1, err := time.Parse(dateFormat, startDate)
	if err != nil {
		return false
	}

	date2, err := time.Parse(dateFormat, endDate)
	if err != nil {
		return false
	}

	// Compare the dates
	if date1.After(date2) {
		return false
	} else {
		return true
	}
}
