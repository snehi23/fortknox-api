package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
)

var localCache *cache.Cache
var authorityMap map[string]bool

func main() {

	localCache = cache.New(5*time.Second, 10*time.Second)
	authorityMap = make(map[string]bool)
	hydrateAuthorityMap(authorityMap)
	handleRequests()
}

func hydrateAuthorityMap(authorityMap map[string]bool) {

	authorityMap["Employee"] = true
	authorityMap["Name"] = true
	authorityMap["Credit_Card"] = true
	authorityMap["Address"] = true
}

func handleRequests() {

	http.Handle("/getToken", http.HandlerFunc(getToken))
	http.Handle("/redeemToken", http.HandlerFunc(redeemToken))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type TokenRequest struct {
	Sde       string
	Authority string
}

type TokenResponse struct {
	Sde       string
	Token     string
	RequestId string
	Authority string
}

type RedeemRequest struct {
	Token     string
	Authority string
}

type RedeemResponse struct {
	Sde       string
	Token     string
	RequestId string
	Authority string
}

func getToken(response http.ResponseWriter, request *http.Request) {

	// fmt.Printf("authority map %v", authorityMap)

	var tokenRequest TokenRequest
	var tokenResponse TokenResponse

	error := json.NewDecoder(request.Body).Decode(&tokenRequest)

	if error != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Somthing is wrong. Please try again later. error: %v sde: %v", error, tokenRequest)
	}

	_, ok := authorityMap[tokenRequest.Authority]

	if !ok {
		log.Printf("Unrecognized authority %v. can not tokenize", tokenRequest.Authority)
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode("Bad Request")
		return
	}

	// fmt.Printf("tokenRequest %+v \n", tokenRequest)

	tokenResponse.Sde = tokenRequest.Sde
	tokenResponse.Token = tokanizeSDE(tokenRequest.Sde)
	token := tokanizeSDE(tokenRequest.Sde)
	localCache.Set(token, tokenRequest.Sde, cache.DefaultExpiration)
	tokenResponse.RequestId = getUUID()

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(tokenResponse)
}

func tokanizeSDE(sde string) string {

	token := b64.StdEncoding.EncodeToString([]byte(sde))
	return token
}

func getUUID() string {

	id := uuid.Must(uuid.NewRandom()).String()
	return id
}

func redeemToken(response http.ResponseWriter, request *http.Request) {

	var redeemRequest RedeemRequest
	var redeemResponse RedeemResponse

	error := json.NewDecoder(request.Body).Decode(&redeemRequest)

	if error != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Somthing is wrong. Please try again later. error: %v token: %v", error, redeemRequest)
	}

	_, ok := authorityMap[redeemRequest.Authority]

	if !ok {
		log.Printf("Unrecognized authority %v. can not redeem", redeemRequest.Authority)
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode("Bad Request")
		return
	}

	// fmt.Printf("redeemRequest %+v \n", redeemRequest)

	sde, found := localCache.Get(redeemRequest.Token)

	// if found in cache, use it
	if found {
		fmt.Printf("redeem from cache")
		redeemResponse.Sde = sde.(string)
		// if not then detokenize it using decoder algo
	} else {
		fmt.Printf("redeem from decoder")
		redeemSDE := detokanizeSDE(redeemRequest.Token)
		redeemResponse.Sde = redeemSDE
	}

	redeemResponse.Token = redeemRequest.Token
	redeemResponse.RequestId = getUUID()

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(redeemResponse)
}

func detokanizeSDE(token string) string {

	sde, _ := b64.StdEncoding.DecodeString(token)
	return string(sde)
}
