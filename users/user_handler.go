package users

import (
	"net/http"
)

type UserHandler struct {
	userService UserServiceIF
}

func NewUserHandler(auth UserServiceIF) *UserHandler {
	return &UserHandler{userService: auth}
}

func (u *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	// TODO 구현할 것!
	panic("not implemented")
}

func (u *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	// TODO 구현할 것!
	panic("not implemented")
}

func (u *UserHandler) GetSelfUser(w http.ResponseWriter, r *http.Request) {
	// TODO 구현할 것!
	panic("not implemented")
}
