package auth

import (
	"github.com/golang-jwt/jwt"
	"time"
)

// TODO 환경변수 등록 후 사용 등 보안적으로 보완해야 함
const secret string = "secret"

func NewClaim(userID string) *jwt.StandardClaims {
	now := time.Now()
	return &jwt.StandardClaims{
		Subject:   userID,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(time.Hour * 24).Unix(),
	}
}

func GenerateToken(claim *jwt.StandardClaims) (string, error) {
	// TODO 구현할 것!
	panic("not implemented")
}

func ValidateToken(token string) (string, error) {
	// TODO 구현할 것!
	panic("not implemented")
}
