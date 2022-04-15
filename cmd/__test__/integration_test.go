package integration_test

import (
	"authentication/auth"
	"authentication/users"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type IntegrationTestSuite struct {
	suite.Suite
	handler *users.UserHandler
}

func (i *IntegrationTestSuite) SetupTest() {
	userRepo := users.NewUserRepository()
	userService := users.NewUserService(userRepo, secret)
	userHandler := users.NewUserHandler(userService)
	i.handler = userHandler
}

func (i *IntegrationTestSuite) TestSuccess() {
	i.T().Run("유저 등록 성공 -> 등록 정보로 토큰 발급 성공 -> 토큰으로 유저 정보 조회 성공", func(t *testing.T) {
		userDto := &users.UserDto{
			UserID:   "testUser",
			Password: "testPassword",
		}
		// 유저 등록 **************************************************************************************
		r, err := createRequest("POST", "/users", userDto, nil)
		if err != nil {
			t.Error(err)
		}
		w := httptest.NewRecorder()
		// 실행
		signUpHandler := http.HandlerFunc(i.handler.SignUp)
		signUpHandler.ServeHTTP(w, r)
		// 검증
		assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
		userResponse, err := getUserResponseFromResponse(w.Result())
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, userDto.UserID, userResponse.UserID)

		// 등록 정보로 토큰 발급 ***************************************************************************
		r, err = createRequest("POST", "/users/token", userDto, nil)
		if err != nil {
			t.Error(err)
		}
		w = httptest.NewRecorder()
		signInHandler := http.HandlerFunc(i.handler.SignIn)
		// 실행
		signInHandler.ServeHTTP(w, r)
		// 검증
		assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
		token, err := getTokenFromResponse(w.Result())
		if err != nil {
			t.Error(err)
		}
		assert.NotNil(t, token)
		assert.NotEmpty(t, token)

		// 토큰으로 유저 정보 조회 ************************************************************************
		headerMap := map[string]string{
			"Authorization": "Bearer " + token,
		}

		r, err = createRequest("GET", "/users/whon_am_i", nil, headerMap)
		if err != nil {
			t.Error(err)
		}
		w = httptest.NewRecorder()
		getSelfUserHandler := http.HandlerFunc(i.handler.GetSelfUser)
		authentication := auth.NewAuthentication(secret)
		stripTokenMiddleWare := authentication.StripTokenMiddleware(getSelfUserHandler)
		// 실행
		stripTokenMiddleWare.ServeHTTP(w, r)
		// 검증
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		userResponse, err = getUserResponseFromResponse(w.Result())
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, userDto.UserID, userResponse.UserID)
	})
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

// 테스트 데이터
const secret string = "testSecret"

func createRequest(method string, url string, userDto *users.UserDto, headers map[string]string) (*http.Request, error) {
	body, err := json.Marshal(userDto)
	if err != nil {
		return nil, err
	}
	request := httptest.NewRequest(method, url, bytes.NewReader(body))
	if len(headers) > 0 {
		for k, v := range headers {
			request.Header.Set(k, v)
		}
	}

	return request, nil
}

func getUserResponseFromResponse(response *http.Response) (*users.UserResponse, error) {
	var userResponse users.UserResponse
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return nil, err
	}
	return &userResponse, nil
}

func getTokenFromResponse(response *http.Response) (string, error) {
	bytesToken, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	var token struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(bytesToken, &token)
	if err != nil {
		return "", err
	}
	return token.Token, nil
}
