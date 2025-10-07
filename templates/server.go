package main

import (
	"html/template"
	"net/http"
)

type Game struct {
	Player1 string
	Player2 string
	Grid    [][]string
	Turn    int // 0 = Player1, 1 = Player2
}

var currentGame *Game

func main() {
	http.HandleFunc("/", startPage)
	http.HandleFunc("/start", startGame)
	http.HandleFunc("/play", playMove)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8080", nil)
}

func startPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/start.html"))
	tmpl.Execute(w, nil)
}

func startGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	player1 := r.FormValue("player1")
	player2 := r.FormValue("player2")

	grid := make([][]string, 6)
	for i := range grid {
		grid[i] = make([]string, 7)
		for j := range grid[i] {
			grid[i][j] = " "
		}
	}

	currentGame = &Game{
		Player1: player1,
		Player2: player2,
		Grid:    grid,
		Turn:    0,
	}

	http.Redirect(w, r, "/play", http.StatusSeeOther)
}

func playMove(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/game.html"))
	tmpl.Execute(w, currentGame)
}
