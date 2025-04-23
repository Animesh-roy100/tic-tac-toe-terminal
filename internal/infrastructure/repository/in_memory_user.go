package repository

import (
	"errors"

	"tic-tac-toe/internal/domain/user"
)

type InMemoryUserRepository struct {
	users map[string]*user.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{users: make(map[string]*user.User)}
}

func (r *InMemoryUserRepository) FindByUsername(username string) (*user.User, error) {
	u, ok := r.users[username]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (r *InMemoryUserRepository) Save(u *user.User) error {
	r.users[u.Username] = u
	return nil
}

func (r *InMemoryUserRepository) All() ([]*user.User, error) {
	var users []*user.User
	for _, u := range r.users {
		users = append(users, u)
	}
	return users, nil
}
