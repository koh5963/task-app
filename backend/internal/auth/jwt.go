package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetUserIDFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", errors.New("authorization header must start with Bearer")
	}

	token := strings.TrimPrefix(authHeader, prefix)
	if token == "" {
		return "", errors.New("token is empty")
	}

	return token, nil
}
