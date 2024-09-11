package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/bson"

	"fortknox-api/repository"
	"fortknox-api/caches"
	"fortknox-api/common"
	"fortknox-api/model"
)

var collection = repository.ConnectToDB().Database("fortknox-db").Collection("tokenStore")
var localCache = caches.SetUpCache()
var authorityMap = common.HydrateAuthorityMap()

func CreateToken(response http.ResponseWriter, request *http.Request) {

	var tokenRequest model.TokenRequest
	var tokenResponse model.TokenResponse

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
		var result model.MongoDBTokenStoreDocument
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
			tokenResponse.Token = common.TokanizeSDE(tokenRequest.Sde)
			localCache.Set(tokenResponse.Token, tokenRequest.Sde, cache.DefaultExpiration)
			localCache.Set(tokenRequest.Sde, tokenResponse.Token, cache.DefaultExpiration)

			newToken := model.MongoDBTokenStoreDocument{Sde: tokenRequest.Sde, Token: tokenResponse.Token}
			_, err := collection.InsertOne(context.TODO(), newToken)

			// if error write to DB then return ERROR
			if err != nil {
				log.Fatalf("Something is wrong while inserting to DB. Please try again later. error: %v sde: %v", err, tokenRequest)
				panic(err)
			}

			log.Printf("token stored in DB")
		}
	}

	tokenResponse.RequestId = common.GetUUID()

	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(tokenResponse)
}

func RedeemToken(response http.ResponseWriter, request *http.Request) {

	var redeemRequest model.RedeemRequest
	var redeemResponse model.RedeemResponse

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

	// if found in cache, return
	if found {
		log.Printf("redeem from cache")
		redeemResponse.Sde = sde.(string)
		// if not found then search in DB
	} else {

		var result model.MongoDBTokenStoreDocument
		filter := bson.D{{Key: "token", Value: redeemRequest.Token}}
		err := collection.FindOne(context.TODO(), filter).Decode(&result)

		// if found, write to cache and return
		if err == nil {
			log.Printf("sde found in DB")
			redeemResponse.Sde = result.Sde
			localCache.Set(redeemRequest.Token, redeemResponse.Sde, cache.DefaultExpiration)
			localCache.Set(redeemResponse.Sde, redeemRequest.Token, cache.DefaultExpiration)
			// if not found, return 404 NOT FOUND
		} else {
			response.WriteHeader(http.StatusNotFound)
			json.NewEncoder(response).Encode("SDE Not Found")
			return
		}

		redeemResponse.Token = redeemRequest.Token
		redeemResponse.RequestId = common.GetUUID()

		response.WriteHeader(http.StatusOK)
		json.NewEncoder(response).Encode(redeemResponse)
	}
}
