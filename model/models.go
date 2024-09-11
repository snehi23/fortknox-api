package model

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