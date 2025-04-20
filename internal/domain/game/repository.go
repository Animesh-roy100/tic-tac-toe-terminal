package game

type GameRepository interface {
	Save(game *Game) error
	Find(id string) (*Game, error)
}
