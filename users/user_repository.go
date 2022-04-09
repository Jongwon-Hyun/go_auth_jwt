package users

type UserRepositoryIF interface {
	Save(user User) *User
	FindByID(id int) (*User, error)
	FindByUserID(userID string) (*User, error)
}
