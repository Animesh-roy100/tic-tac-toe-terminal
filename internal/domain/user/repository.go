package user

type UserRepository interface {
	Save(user *User) error
	Find(username string) (*User, error)
	All() ([]*User, error)
}
