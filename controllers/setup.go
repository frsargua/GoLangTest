package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func New() http.Handler {
  router := mux.NewRouter()


  router.HandleFunc("/spots", GetSpots).Methods("GET")

  
  
  return router
}