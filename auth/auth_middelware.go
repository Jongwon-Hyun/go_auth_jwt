package auth

import (
	"net/http"
)

func AuthenticateMiddleware(next http.Handler) http.Handler {
	// TODO 구현할 것!
	panic("not implemented")
}

const (
	AuthSchema string = "Bearer "
)

func getTokenFromRequest(r *http.Request) (string, error) {
	// TODO 구현할 것!
	panic("not implemented")
}
