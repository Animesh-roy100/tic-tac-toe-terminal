package repository

import (
	"errors"
	"tic-tac-toe/internal/domain/game"
)

type InMemoryGameRepository struct {
	games map[string]*game.Game
}

func NewInMemoryGameRepository() *InMemoryGameRepository {
	return &InMemoryGameRepository{games: make(map[string]*game.Game)}
}

func (r *InMemoryGameRepository) FindByID(id string) (*game.Game, error) {
	g, ok := r.games[id]
	if !ok {
		return nil, errors.New("game not found")
	}
	return g, nil
}

func (r *InMemoryGameRepository) Save(g *game.Game) error {
	r.games[g.ID] = g
	return nil
}
