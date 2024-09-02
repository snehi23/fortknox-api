package main

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/bson"
)

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

type CacheRequest struct {
	Key   string
	Value string
}

type MongoDBTokenStoreDocument struct {
	Sde   string `json:"sde"`
	Token string `json:"token"`
}

var collection = connectToDB().Database("fortknox-db").Collection("tokenStore")
var localCache = setUpCache()

func createToken(response http.ResponseWriter, request *http.Request) {

	var tokenRequest TokenRequest
	var tokenResponse TokenResponse

	error := json.NewDecoder(request.Body).Decode(&tokenRequest)

	if error != nil {
		response.WriteHeader(http.StatusInternalServerError)
		log.Fatalf("Something is wrong while decoding the request. Please try again later. error: %v sde: %v", error, tokenRequest)
	}

	_, ok := authorityMap[tokenRequest.Authority]

	if !ok {
		log.Printf("Unrecognized authority %v. can not tokenize", tokenRequest.Authority)
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode("Bad Request")
		return
	}

	tokenResponse.Authority = tokenRequest.Authority
	tokenResponse.Sde = tokenRequest.Sde

	// search in cache
	token, found := localCache.Get(tokenRequest.Sde)

	// if found, return
	if found {
		log.Printf("token found in cache")
		tokenResponse.Token = token.(string)
	// if not found, search in DB	
	} else {
		var result MongoDBTokenStoreDocument

		filter := bson.D{{Key: "sde", Value: tokenRequest.Sde}}

		err := collection.FindOne(context.TODO(), filter).Decode(&result)

		// if found, write to cache and return
		if err == nil {

			log.Printf("token found in DB")
			tokenResponse.Token = result.Sde
			localCache.Set(tokenResponse.Token, tokenRequest.Sde, cache.DefaultExpiration)
			localCache.Set(tokenRequest.Sde, tokenResponse.Token, cache.DefaultExpiration)
		// if not found, generate new token, write to DB, write to cache and return	
		} else {
			tokenResponse.Token = tokanizeSDE(tokenRequest.Sde)
			localCache.Set(tokenResponse.Token, tokenRequest.Sde, cache.DefaultExpiration)
			localCache.Set(tokenRequest.Sde, tokenResponse.Token, cache.DefaultExpiration)

			newToken := MongoDBTokenStoreDocument{Sde: tokenRequest.Sde, Token: tokenResponse.Token}
			_, err := collection.InsertOne(context.TODO(), newToken)

			// if error write to DB then return ERROR
			if err != nil {
				log.Fatalf("Something is wrong while inserting to DB. Please try again later. error: %v sde: %v", err, tokenRequest)
				panic(err)
			}

			log.Printf("token stored in DB")
		}
	}

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

	redeemResponse.Authority = redeemRequest.Authority

	sde, found := localCache.Get(redeemRequest.Token)

	// if found in cache, use it
	if found {
		log.Printf("redeem from cache")
		redeemResponse.Sde = sde.(string)
	// if not then detokenize it using decoder algo
	} else {
		log.Printf("redeem from decoder")
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
