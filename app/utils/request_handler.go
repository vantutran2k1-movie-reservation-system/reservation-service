package utils

import (
	"net/http"
	"strings"
)

func GetAuthorizationHeader(req *http.Request) string {
	token := req.Header.Get("Authorization")
	tokenParts := strings.Split(token, " ")
	if len(tokenParts) == 2 {
		token = tokenParts[1]
	}

	return token
}
