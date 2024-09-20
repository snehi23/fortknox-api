package common

import (
	b64 "encoding/base64"
	"net/http"
	"os"

	"github.com/google/uuid"
)

func HydrateAuthorityMap() map[string]bool {

	var authorityMap = make(map[string]bool)

	authorityMap["Employee"] = true
	authorityMap["Name"] = true
	authorityMap["Credit_Card"] = true
	authorityMap["Address"] = true

	return authorityMap
}

func TokanizeSDE(sde string) string {

	token := b64.StdEncoding.EncodeToString([]byte(sde))
	return token
}

func GetUUID() string {

	id := uuid.Must(uuid.NewRandom()).String()
	return id
}

func DetokanizeSDE(token string) string {

	sde, _ := b64.StdEncoding.DecodeString(token)
	return string(sde)
}

func GetAPIKey() string {

	apiKey := os.Getenv("FORTKNOX_API_KEY")
	return apiKey
}

func IsValidAPIKey(request *http.Request) bool {
	// Check if the API key is provided in the request header
	providedKey := request.Header.Get("X-API-Key")
	apiKey := GetAPIKey()

	return providedKey == apiKey
}
