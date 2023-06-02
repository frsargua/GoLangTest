package controllers

import (
	"math"
	"regexp"
)

const (
	earthRadius = 6356752 
	pi          = math.Pi
)


func ConvertMetersToCoordinates(distance float64, latitude, longitude float64) (float64, float64) {
	latRad := latitude * pi / 180
	lonRad := longitude * pi / 180

	angularDistance := distance / earthRadius

	// Calculate the new latitude in radians
	newLatRad := latRad + angularDistance
	// Calculate the new latitude in degrees
	newLatitude := newLatRad * 180 / pi

	newLonRad := lonRad + angularDistance/math.Cos(latRad)
	newLongitude := newLonRad * 180 / pi

	return newLatitude, newLongitude
}



func isValidCoordinate(coordinate string) bool {
	regex := `^[-+]?([1-8]?\d(\.\d+)?|90(\.0+)?)$`
	match, _ := regexp.MatchString(regex, coordinate)
	return match
}