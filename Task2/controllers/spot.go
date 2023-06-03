package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/frsargua/GoLangTest/models"
)


func GetSpots(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	queryParams := r.URL.Query()

	
	latitude, longitude, radius, isCircle, err := verifyAndParseParameters(queryParams)

	spots, err := getSpotsAroundCoordinate(latitude, longitude, radius, isCircle)


	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(spots)
  }

func getSpotsAroundCoordinate(latitude, longitude, radius float64, isCircle bool) ([]models.Spot, error) {
	var spots []models.Spot

	query,err := getQueryString(latitude, longitude, radius, isCircle)
	if err != nil {
		return spots, err
	}

	err2 := models.DB.Raw(query).Scan(&spots).Error

	if err2 != nil {
		return spots, errors.New("error retrieving data")
	}


	return spots, nil
}

func getQueryString(latitude, longitude, radius float64, isCircle bool) (string,error) {
	query := ""

	newLatitude, newLongitude  := ConvertMetersToCoordinates(radius,latitude,longitude)
	newLatitudeNeg, newLongitudeNeg  := ConvertMetersToCoordinates(-radius,latitude,longitude)


	if(isCircle){
		query = fmt.Sprintf(`
		WITH clusters AS (
		  SELECT *, ST_ClusterDBSCAN(coordinates::geometry, eps := 0.00045, minpoints := 2) OVER () AS cluster_id
		  FROM public."MY_TABLE"
		  WHERE ST_DWithin(coordinates::geography, ST_SetSRID(ST_MakePoint(%f, %f), 4326), %f)
		  ), 
		  centroids AS (
			SELECT cluster_id,
				   ST_Centroid(ST_Collect(coordinates::geometry)) AS cluster_centroid
			FROM clusters
			WHERE cluster_id IS NOT NULL
			GROUP BY cluster_id
		)
		SELECT clusters.id, clusters.name, clusters.website, clusters.coordinates, clusters.description, clusters.rating
		FROM clusters
		LEFT JOIN centroids 
		ON clusters.cluster_id = centroids.cluster_id
		ORDER BY 
			CASE 
				WHEN clusters.cluster_id IS NULL THEN ST_Distance(clusters.coordinates::geometry, ST_SetSRID(ST_MakePoint(%f, %f), 4326))
				ELSE ST_Distance(centroids.cluster_centroid::geometry, ST_SetSRID(ST_MakePoint(%f, %f), 4326))
			END,
			rating
		  
	`, longitude,latitude, radius, longitude,latitude,longitude,latitude)
		}else{
			query = fmt.Sprintf(`
			WITH clusters AS (
				SELECT *,
					   ST_ClusterDBSCAN(coordinates::geometry, eps := 0.00045, minpoints := 2) OVER () AS cluster_id
				FROM public."MY_TABLE"
				WHERE ST_Intersects(
					ST_MakeEnvelope(
						%f, %f,
						%f, %f,
						4326
					),
					coordinates::geometry
				)
			), 
			centroids AS (
				SELECT cluster_id,
					   ST_Centroid(ST_Collect(coordinates::geometry)) AS cluster_centroid
				FROM clusters
				WHERE cluster_id IS NOT NULL
				GROUP BY cluster_id
			)
			SELECT clusters.id, clusters.name, clusters.website, clusters.coordinates, clusters.description, clusters.rating
			FROM clusters
			LEFT JOIN centroids 
			ON clusters.cluster_id = centroids.cluster_id
			ORDER BY 
				CASE 
					WHEN clusters.cluster_id IS NULL THEN ST_Distance(clusters.coordinates::geometry, ST_SetSRID(ST_MakePoint(%f, %f), 4326))
					ELSE ST_Distance(centroids.cluster_centroid::geometry, ST_SetSRID(ST_MakePoint(%f, %f), 4326))
				END,
				rating
			
			`, newLongitudeNeg, newLatitudeNeg, newLongitude, newLatitude, longitude,latitude,longitude,latitude)
	
	}

	if query == "" {
		return query, errors.New("error creating query string.")
	}
		return query, nil
}

func verifyAndParseParameters(queryParams url.Values) (float64, float64, float64, bool, error) {
	latitudeStr, longitudeStr, radiusStr, isCircleStr :=
		queryParams.Get("latitude"), queryParams.Get("longitude"), queryParams.Get("radius"), queryParams.Get("isCircle")

	if latitudeStr == "" || longitudeStr == "" || radiusStr == "" || isCircleStr == "" {
		return 0, 0, 0, false, errors.New("one or more parameters are missing")
	}

	if !isValidCoordinate(latitudeStr) || !isValidCoordinate(longitudeStr) {
		return 0, 0, 0, false, errors.New("invalid latitude or longitude")
	}

	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil {
		return 0, 0, 0, false, errors.New("invalid radius")
	}

	isCircle, err := strconv.ParseBool(isCircleStr)
	if err != nil {
		return 0, 0, 0, false, errors.New("invalid isCircle")
	}

	latitude, _ := strconv.ParseFloat(latitudeStr, 64)
	longitude, _ := strconv.ParseFloat(longitudeStr, 64)

	return latitude, longitude, radius, isCircle, nil
}