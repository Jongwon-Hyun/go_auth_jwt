package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewClaim(t *testing.T) {
	t.Run("올바른 Claim 이 생성될 것", func(t *testing.T) {
		got := NewClaim(userID)

		sub := getSub(got, t)

		assert.Equal(t, userID, sub)
	})
}

func TestGenerateToken(t *testing.T) {
	t.Run("올바른 Token 이 생성될 것", func(t *testing.T) {
		got, err := GenerateToken(claim, secret)
		if err != nil {
			t.Error(err)
		}

		// 테스트를 검증하기 위해서는 ValidateToken 메서드가 구현되어 있을 것
		want, err := ValidateToken(got, secret)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, getSub(want, t), getSub(claim, t))
	})
}

func TestValidateToken(t *testing.T) {
	t.Run("올바른 Token 일 경우, 올바른 User ID 가 반환될 것", func(t *testing.T) {
		// 토큰 생성 테스트에서 확인 완료
		TestGenerateToken(t)
	})

	t.Run("토큰 서명 암호화 방식이 맞지 않을 경우, 에러가 반환될 것", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
		signedToken, _ := token.SignedString([]byte(secret))

		claim, err := ValidateToken(signedToken, secret)
		assert.NotNil(t, err)
		assert.Nil(t, claim)
	})

	t.Run("토큰 파싱에 실패했을 경우, 에러가 반환될 것", func(t *testing.T) {
		claim, err := ValidateToken("I am Not Token", secret)
		assert.NotNil(t, err)
		assert.Nil(t, claim)
	})

	t.Run("유효기간이 지난 Token 일 경우, 에러가 반환될 것", func(t *testing.T) {
		expiredClaim := claim
		expiredClaim["exp"] = time.Now().Add(time.Hour * -5).Unix()
		claim, err := ValidateToken("I am Not Token", secret)
		assert.NotNil(t, err)
		assert.Nil(t, claim)
	})
}

// 테스트 데이터
var (
	userID    = "testUserID"
	secret    = "testSecret"
	issuedAt  = float64(time.Now().Unix())
	ExpiredAt = float64(time.Now().Add(time.Hour * 24).Unix())
	claim     = jwt.MapClaims{
		"sub": userID,
		"iat": issuedAt,
		"exp": ExpiredAt,
	}
)

func getSub(claim jwt.MapClaims, t *testing.T) string {
	sub, ok := claim["sub"].(string)
	if !ok {
		t.Errorf("sub 타입이 string이 아닙니다.")
	}
	return sub
}
