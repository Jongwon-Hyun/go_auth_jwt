package auth

import (
	"authentication"
	"context"
	"errors"
	"net/http"
	"strings"
)

type authenticationMiddleware struct {
	secret string
}

func NewAuthentication(secret string) *authenticationMiddleware {
	return &authenticationMiddleware{secret: secret}
}

func (a *authenticationMiddleware) StripTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getTokenFromRequest(r)
		if err != nil {
			http.Error(w, err.Error(), authentication.ErrStatusCode(err))
			return
		}

		userID, err := ValidateToken(token, a.secret)
		if err != nil {
			http.Error(w, err.Error(), authentication.ErrStatusCode(err))
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

const (
	BearerSchema string = "BEARER "
)

func getTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New(authentication.ErrAuthorizationHeaderRequired)
	}

	bearerLength := len(BearerSchema)
	if len(authHeader) > bearerLength && strings.ToUpper(authHeader[0:bearerLength]) == BearerSchema {
		return authHeader[bearerLength:], nil
	}

	return "", errors.New(authentication.ErrInvalidBearerScheme)
}
