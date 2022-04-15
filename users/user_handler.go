package users

import (
	"authentication"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserHandler struct {
	userService UserServiceIF
}

func NewUserHandler(userService UserServiceIF) *UserHandler {
	return &UserHandler{userService: userService}
}

func (u *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var userDto UserDto
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		http.Error(w, err.Error(), authentication.ErrStatusCode(err))
		return
	}

	userResponse, err := u.userService.SignUp(userDto)
	if err != nil {
		http.Error(w, err.Error(), authentication.ErrStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(userResponse)
	if err != nil {
		http.Error(w, err.Error(), authentication.ErrStatusCode(err))
	}
}

func (u *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var userDto UserDto
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		http.Error(w, err.Error(), authentication.ErrStatusCode(err))
		return
	}

	token, err := u.userService.SignIn(userDto)
	if err != nil {
		http.Error(w, err.Error(), authentication.ErrStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, err = fmt.Fprintf(w, "{\"token\":\"%s\"}", token)
	if err != nil {
		http.Error(w, err.Error(), authentication.ErrStatusCode(err))
	}
}

func (u *UserHandler) GetSelfUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("user_id").(string)
	if !ok {
		http.Error(w, "unauthorized user", http.StatusUnauthorized)
		return
	}

	userResponse, err := u.userService.GetUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), authentication.ErrStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(userResponse)
	if err != nil {
		http.Error(w, err.Error(), authentication.ErrStatusCode(err))
	}
}
