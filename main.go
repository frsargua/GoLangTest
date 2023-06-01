package main

import (
	"net/http"

	"github.com/frsargua/GoLangTest/controllers"
	"github.com/frsargua/GoLangTest/models"
	"github.com/joho/godotenv"
)


func main() {
  godotenv.Load()

  handler := controllers.New() 

  server := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: handler,
  }



  models.ConnectDatabase()

  server.ListenAndServe()
}