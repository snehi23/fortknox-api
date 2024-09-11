package common

import (
	b64 "encoding/base64"

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
