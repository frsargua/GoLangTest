package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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
	if !isValidCoordinate(latitude) || ! isValidCoordinate(longitude) {
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

lat, nil := strconv.ParseFloat(latitude,64)
lon, nil := strconv.ParseFloat(longitude,64)
rad := float64(radius)
	rows, err := getSpotsInArea(lat,lon,rad,isCircle);
  
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(rows)
  }


func getSpotsInArea(latitude float64, longitude float64, radius float64, isCircle bool) ([]models.Spot, error) {
	query := ""

	newLatitude, newLongitude  := ConvertMetersToCoordinates(radius,latitude,longitude)
	newLatitudeNeg, newLongitudeNeg  := ConvertMetersToCoordinates(-radius,latitude,longitude)


	if(isCircle){
		query = fmt.Sprintf(`
		  SELECT *
		  FROM (
			  SELECT 
				  *, ST_X(ST_Transform(t1.coordinates, 4326)) AS longitude, ST_Y(ST_Transform(t1.coordinates, 4326)) AS latitude
			  FROM public."MY_TABLE" AS t1
		  ) AS TB 
		  WHERE ST_DWithin(coordinates::geography, ST_SetSRID(ST_MakePoint(%f, %f), 4326), %f)
		  ORDER BY
			  CASE
				  WHEN ST_Distance(coordinates::geography, ST_SetSRID(ST_MakePoint(TB.longitude, TB.latitude), 4326)) < 50 THEN rating
				  ELSE NULL
			  END,
			  ST_Distance(coordinates::geography, ST_SetSRID(ST_MakePoint(%f, %f), 4326));
		  
	`, longitude,latitude, radius, longitude,latitude)
		}else{
			query = fmt.Sprintf(`
			SELECT *
			FROM (
				SELECT 
					*, ST_X(ST_Transform(t1.coordinates, 4326)) AS longitude, ST_Y(ST_Transform(t1.coordinates, 4326)) AS latitude
				FROM public."MY_TABLE" AS t1
			) AS TB 
			WHERE ST_Intersects(
				ST_MakeEnvelope(
					%f, %f,
					%f, %f,
					4326
				),
				coordinates::geometry
			)
			ORDER BY
		  CASE
			WHEN ST_Distance(coordinates::geography, ST_SetSRID(ST_MakePoint(TB.longitude, TB.latitude), 4326)) < 50 THEN rating
			ELSE NULL
		  END,
		  ST_Distance(coordinates::geography, ST_SetSRID(ST_MakePoint(%f, %f), 4326));
			`, newLongitudeNeg, newLatitudeNeg, newLongitude, newLatitude, longitude,latitude)
	
	}

	var spots []models.Spot
	err := models.DB.Raw(query).Scan(&spots).Error
	
	if err != nil {
		err = errors.New("Error retrieving data")
		return spots, err
	}
	


	return spots, nil
}