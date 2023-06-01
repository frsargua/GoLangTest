package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/frsargua/GoLangTest/models"
)


func GetSpots(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
  
	var spots []models.Spot
	models.DB.Find(&spots)
  
	json.NewEncoder(w).Encode(spots)
  }