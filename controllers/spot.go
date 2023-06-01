package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/frsargua/GoLangTest/models"
)


func GetSpots(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	queryParams := r.URL.Query()
	latitude := queryParams.Get("latitude")
	longitude := queryParams.Get("longitude") 
	radiusStr := queryParams.Get("radius") 
	isCircleStr := queryParams.Get("isCircle") 


		// Check if any parameter is missing
		if latitude == "" || longitude == "" || radiusStr == "" || isCircleStr == "" {
			err := errors.New("One or more parameters are missing")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

			// Verify latitude and longitude format
	if !isValidCoordinate(latitude) || !isValidCoordinate(longitude) {
		err := errors.New("Invalid latitude or longitude")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Verify radius is an integer
	radius, err := strconv.Atoi(radiusStr)
	if err != nil {
		err := errors.New("Invalid radius")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Verify isCircle is a boolean
	isCircle, err := strconv.ParseBool(isCircleStr)
	if err != nil {
		err := errors.New("Invalid isCircle")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// var results []models.Spot
	
	rows, err := getSpotsInArea(latitude,longitude,radius,isCircle);
  
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(rows)
  }

  func isValidCoordinate(coordinate string) bool {
	regex := `^[-+]?([1-8]?\d(\.\d+)?|90(\.0+)?)$`
	match, _ := regexp.MatchString(regex, coordinate)
	return match
}

func getSpotsInArea(latitude string, longitude string, radius int, isCircle bool) ([]models.Spot, error) {
	var spots []models.Spot
	var query string
	if(isCircle){
		query = fmt.Sprintf(`SELECT *
		FROM public."MY_TABLE"
		WHERE ST_DWithin(coordinates::geography, ST_SetSRID(ST_MakePoint(%s, %s), 4326), %d)
		ORDER BY
		  CASE
			WHEN ST_Distance(coordinates::geography, ST_SetSRID(ST_MakePoint(%s, %s), 4326)) < 50 THEN rating
			ELSE NULL
		  END,
		  ST_Distance(coordinates::geography, ST_SetSRID(ST_MakePoint(%s,%s), 4326));
		`,latitude,longitude,radius,latitude,longitude,latitude,longitude)
		}else{
		query = fmt.Sprintf(`SELECT * FROM public."MY_TABLE" LIMIT 10`)
	}


	err := models.DB.Limit(10).Raw(query).Scan(&spots).Error

	fmt.Println(latitude , longitude , radius, isCircle)

	if err != nil {
		err := errors.New("Error retrieving data")
		return spots, err
	}


	return spots, nil
}