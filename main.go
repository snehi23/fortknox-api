package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var authorityMap map[string]bool

func main() {
	authorityMap = make(map[string]bool)
	hydrateAuthorityMap(authorityMap)
	route := mux.NewRouter()
	// Base Path
	s := route.PathPrefix("/fortknox").Subrouter()
	// Routes
	s.HandleFunc("/createToken", createToken).Methods("POST")
	s.HandleFunc("/redeemToken", redeemToken).Methods("GET")
	// Run Server
	log.Fatal(http.ListenAndServe(":8080", s))
}

func hydrateAuthorityMap(authorityMap map[string]bool) {

	authorityMap["Employee"] = true
	authorityMap["Name"] = true
	authorityMap["Credit_Card"] = true
	authorityMap["Address"] = true
}
