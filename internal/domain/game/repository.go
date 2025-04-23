package game

type GameRepository interface {
	Save(game *Game) error
	FindByID(id string) (*Game, error)
}
