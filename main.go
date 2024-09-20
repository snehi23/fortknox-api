package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"fortknox-api/common"
	"fortknox-api/service"
)

func main() {

	// note : how-to : setup API Key
	// export FORTKNOX_API_KEY="<api-key>"
	apiKey := common.GetAPIKey()
	if apiKey == "" {
		log.Printf("apikey is not set for app.")
		return
	}

	route := mux.NewRouter()
	// Base Path
	s := route.PathPrefix("/fortknox").Subrouter()
	// Routes
	s.HandleFunc("/createToken", service.CreateToken).Methods("POST")
	s.HandleFunc("/redeemToken", service.RedeemToken).Methods("GET")
	log.Printf("Server is running on http://localhost:8080")
	// Run Server
	log.Fatal(http.ListenAndServe(":8080", s))
}
