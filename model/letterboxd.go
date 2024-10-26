package model

type Letterboxd struct {
	Email     string   `db:"email"`
	Username  string   `db:"username"`
	Watchlist []string `db:"watchlist"`
}
