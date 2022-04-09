package authentication

import "net/http"

const (
	ErrUserNotFound                = "user not found"
	ErrUserAlreadyExists           = "user already exists"
	ErrInvalidPassword             = "invalid password"
	ErrTokenGenerationFailed       = "fail generate token"
	ErrAuthorizationHeaderRequired = "authorization header required"
	ErrInvalidBearerScheme         = "invalid Bearer scheme"
	ErrInvalidToken                = "invalid token"
)

func ErrStatusCode(err error) int {
	switch err.Error() {
	case ErrUserNotFound:
		return http.StatusNotFound
	case ErrUserAlreadyExists:
		return http.StatusBadRequest
	case ErrInvalidPassword:
		return http.StatusBadRequest
	case ErrTokenGenerationFailed:
		return http.StatusInternalServerError
	case ErrAuthorizationHeaderRequired:
		return http.StatusBadRequest
	case ErrInvalidBearerScheme:
		return http.StatusBadRequest
	case ErrInvalidToken:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
