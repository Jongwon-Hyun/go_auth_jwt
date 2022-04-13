package domain

type User struct {
	ID       int
	UserID   string
	Password string
}

type UserWithoutPassword struct {
	ID     int
	UserID string
}
