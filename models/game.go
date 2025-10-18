package models

type Game struct {
	Player1      string
	Player2      string
	Grid         [][]string
	Turn         int
	Winner       string
	Difficulty   string
	Player1Score int
	Player2Score int
}

var CurrentGame *Game
