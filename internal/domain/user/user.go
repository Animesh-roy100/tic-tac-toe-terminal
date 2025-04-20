package user

type User struct {
	Username  string
	Score     int
	WinStreak int
}

func NewUser(username string) *User {
	return &User{
		Username: username,
	}
}

func (u *User) WinGame(isAIGame bool) string {
	return ""
}

func (u *User) LoseGame() {
	u.WinStreak = 0
}

func (u *User) DrawGame() {
	u.WinStreak = 0
}
