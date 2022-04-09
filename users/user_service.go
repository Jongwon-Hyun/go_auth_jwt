package users

type UserServiceIF interface {
	SignUp(userDto UserDto) (UserResponse, error)
	SignIn(userDto UserDto) (string, error)
	GetUserByUserID(userID string) (UserResponse, error)
}
