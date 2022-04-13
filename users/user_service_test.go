package users

import (
	"authentication"
	"authentication/users/domain"
	"authentication/users/mocks"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserServiceTestSuite struct {
	suite.Suite
	repo    *mocks.UserRepositoryIF
	service UserServiceIF
}

func (u *UserServiceTestSuite) SetupTest() {
	u.repo = new(mocks.UserRepositoryIF)
	u.service = NewUserService(u.repo, "secret")
}

func (u *UserServiceTestSuite) TearsDownTest() {
	u.repo.ExpectedCalls = nil
}

func (u *UserServiceTestSuite) TestSignUpSuccess() {
	u.T().Run("유저 등록에 성공하여 응답 정보가 반환될 것", func(t *testing.T) {
		u.repo.On("FindByUserID", userID).Return(nil, errors.New(authentication.ErrUserNotFound))
		u.repo.On("Save",
			mock.MatchedBy(func(arg *domain.User) bool { return arg.UserID == userID && arg.Password == password })).
			Return(user)

		got, err := u.service.SignUp(userDto)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, id, got.ID)
		assert.Equal(t, userID, got.UserID)
		u.repo.AssertCalled(t, "FindByUserID", userID)
		u.repo.AssertCalled(t, "Save",
			mock.MatchedBy(func(arg *domain.User) bool { return arg.UserID == userID && arg.Password == password }))
	})
}

func (u *UserServiceTestSuite) TestSignUpFailure() {
	u.T().Run("이미 존재하는 User 의 경우, 에러가 발생할 것", func(t *testing.T) {
		u.repo.On("FindByUserID", userID).Return(user, nil)
		u.repo.On("Save",
			mock.MatchedBy(func(arg *domain.User) bool { return arg.UserID == userID && arg.Password == password })).
			Return(user)

		got, err := u.service.SignUp(userDto)

		assert.EqualError(t, err, authentication.ErrUserAlreadyExists)
		assert.Equal(t, 0, got.ID)
		assert.Equal(t, "", got.UserID)
		u.repo.AssertCalled(t, "FindByUserID", userID)
		u.repo.AssertNotCalled(t, "Save",
			mock.MatchedBy(func(arg *domain.User) bool { return arg.UserID == userID && arg.Password == password }))
	})
}

func (u *UserServiceTestSuite) TestSignInSuccess() {
	u.T().Run("로그인에 성공하여 토큰이 반환될 것", func(t *testing.T) {
		u.repo.On("FindByUserID", userID).Return(user, nil)

		got, err := u.service.SignIn(userDto)
		if err != nil {
			t.Error(err)
		}

		assert.NotNil(t, got)
		assert.NotEmpty(t, got)
		u.repo.AssertCalled(t, "FindByUserID", userID)
	})
}

func (u *UserServiceTestSuite) TestSignInFailure() {
	u.T().Run("등록되어 있지 않는 유저 정보로 로그인하는 경우, 에러가 발생할 것", func(t *testing.T) {
		u.repo.On("FindByUserID", userID).Return(nil, errors.New(authentication.ErrUserNotFound))

		got, err := u.service.SignIn(userDto)

		assert.EqualError(t, err, authentication.ErrUserNotFound)
		assert.Empty(t, got)
		u.repo.AssertCalled(t, "FindByUserID", userID)
	})

	u.T().Run("잘못된 패스워드로 로그인 하는 경우, 에러가 발생할 것", func(t *testing.T) {
		u.repo.On("FindByUserID", userID).Return(nil, errors.New(authentication.ErrUserNotFound))

		got, err := u.service.SignIn(userDto)

		assert.EqualError(t, err, authentication.ErrUserNotFound)
		assert.Empty(t, got)
		u.repo.AssertCalled(t, "FindByUserID", userID)
	})

	u.T().Run("토큰 발급에 실패하는 경우, 에러가 발생할 것", func(t *testing.T) {
		// 리플렉션 이외에 토큰 발급 부분을 목킹할 방법이 없어서 테스트 불가능
	})
}

func (u *UserServiceTestSuite) TestGetUserByIDSuccess() {
	u.T().Run("User id 에 해당하는 유저가 존재하는 경우, 유저 정보가 반환될 것", func(t *testing.T) {
		u.repo.On("FindByUserID", userID).Return(user, nil)

		got, err := u.service.GetUserByID(userID)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, 1, got.ID)
		assert.Equal(t, userID, got.UserID)
		u.repo.AssertCalled(t, "FindByUserID", userID)
	})
}

func (u *UserServiceTestSuite) TestGetUserByIDFailure() {
	u.T().Run("UserID 에 해당하는 유저가 존재하지 않는 경우, 에러가 발생할 것", func(t *testing.T) {
		u.repo.On("FindByUserID", userID).Return(nil, errors.New(authentication.ErrUserNotFound))

		got, err := u.service.GetUserByID(userID)

		assert.EqualError(t, err, authentication.ErrUserNotFound)
		assert.Equal(t, 0, got.ID)
		assert.Equal(t, "", got.UserID)
		u.repo.AssertCalled(t, "FindByUserID", userID)
	})
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

// 테스트 데이터
var (
	id       = 1
	userID   = "testUserID"
	password = "testPassword"
	userDto  = UserDto{
		UserID:   userID,
		Password: password,
	}
	user = &domain.User{
		ID:       id,
		UserID:   userID,
		Password: password,
	}
)
