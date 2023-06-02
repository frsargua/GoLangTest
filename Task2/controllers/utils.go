package controllers

import "regexp"

func isValidCoordinate(coordinate string) bool {
	regex := `^[-+]?([1-8]?\d(\.\d+)?|90(\.0+)?)$`
	match, _ := regexp.MatchString(regex, coordinate)
	return match
}