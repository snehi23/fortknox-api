package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"fortknox-api/service"
)

func main() {
	route := mux.NewRouter()
	// Base Path
	s := route.PathPrefix("/fortknox").Subrouter()
	// Routes
	s.HandleFunc("/createToken", service.CreateToken).Methods("POST")
	s.HandleFunc("/redeemToken", service.RedeemToken).Methods("GET")
	// Run Server
	log.Fatal(http.ListenAndServe(":8080", s))
}
