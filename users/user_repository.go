package users

import (
	"authentication"
	"authentication/users/domain"
	"errors"
)

type UserRepositoryIF interface {
	Save(user *domain.User) *domain.User
	FindByID(id int) (*domain.User, error)
	FindByUserID(userID string) (*domain.User, error)
}

type userRepository struct {
	ID    int
	Users map[int]*domain.User
}

func NewUserRepository() UserRepositoryIF {
	return &userRepository{
		ID:    0,
		Users: make(map[int]*domain.User),
	}
}

func (u *userRepository) Save(user *domain.User) *domain.User {
	u.ID++
	user.ID = u.ID
	u.Users[u.ID] = user

	return user
}

func (u *userRepository) FindByID(id int) (*domain.User, error) {
	user := u.Users[id]
	if user == nil {
		return nil, errors.New(authentication.ErrUserNotFound)
	}
	return user, nil
}

func (u *userRepository) FindByUserID(userID string) (*domain.User, error) {
	for _, user := range u.Users {
		if user.UserID == userID {
			return user, nil
		}
	}
	return nil, errors.New(authentication.ErrUserNotFound)
}
