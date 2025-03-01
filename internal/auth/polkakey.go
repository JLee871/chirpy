package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	token := headers.Get("Authorization")
	if token == "" {
		return token, errors.New("no auth header")
	}

	split := strings.Split(token, " ")
	if len(split) < 2 || split[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}

	return split[1], nil
}
