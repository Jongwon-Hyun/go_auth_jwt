package users

import (
	"authentication"
	"authentication/auth"
	"authentication/users/domain"
	"errors"
)

type UserServiceIF interface {
	SignUp(userDto UserDto) (UserResponse, error)
	SignIn(userDto UserDto) (string, error)
	GetUserByID(userID string) (UserResponse, error)
}

type userService struct {
	userRepository UserRepositoryIF
	secret         string
}

func NewUserService(userRepository UserRepositoryIF, secret string) UserServiceIF {
	return &userService{
		userRepository: userRepository, secret: secret,
	}
}

func (u *userService) SignUp(userDto UserDto) (UserResponse, error) {
	user := &domain.User{
		UserID:   userDto.UserID,
		Password: userDto.Password,
	}
	foundUser, _ := u.userRepository.FindByUserID(user.UserID)
	if foundUser != nil {
		return UserResponse{}, errors.New(authentication.ErrUserAlreadyExists)
	}

	persistedUser := u.userRepository.Save(user)

	return UserResponse{
		ID:     persistedUser.ID,
		UserID: persistedUser.UserID,
	}, nil
}

func (u *userService) SignIn(userDto UserDto) (string, error) {
	user, err := u.userRepository.FindByUserID(userDto.UserID)
	if err != nil {
		return "", err
	}

	if user.Password != userDto.Password {
		return "", errors.New(authentication.ErrInvalidPassword)
	}

	token, err := auth.GenerateToken(auth.NewClaim(user.UserID), u.secret)
	if err != nil {
		return "", errors.New(authentication.ErrTokenGenerationFailed)
	}
	return token, nil
}

func (u *userService) GetUserByID(userID string) (UserResponse, error) {
	user, err := u.userRepository.FindByUserID(userID)
	if err != nil {
		return UserResponse{}, err
	}
	return UserResponse{
		ID:     user.ID,
		UserID: user.UserID,
	}, nil
}
