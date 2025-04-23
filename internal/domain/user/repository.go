package user

type UserRepository interface {
	Save(user *User) error
	FindByUsername(username string) (*User, error)
	All() ([]*User, error)
}
