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
	if isAIGame {
		u.Score += 1
	} else {
		u.Score += 2
	}
	u.WinStreak++
	bonusMsg := ""
	if u.WinStreak == 3 {
		u.Score += 5
		bonusMsg = "You earned 5 bonus points for a 3-game win streak!"
	} else if u.WinStreak == 5 {
		u.Score += 10
		bonusMsg = "You earned 10 bonus points for a 5-game win streak!"
	}
	return bonusMsg
}

func (u *User) LoseGame() {
	u.WinStreak = 0
}

func (u *User) DrawGame() {
	u.WinStreak = 0
}
