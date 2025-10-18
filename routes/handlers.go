package routes

import (
	"html/template"
	"net/http"
	"power4/models"
	"strconv"
)

func StartPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/start.html"))
	tmpl.Execute(w, nil)
}

func StartGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	player1 := r.FormValue("player1")
	player2 := r.FormValue("player2")
	difficulty := r.FormValue("difficulty")

	rows, cols := 6, 7
	switch difficulty {
	case "facile":
		rows, cols = 6, 7
	case "moyen":
		rows, cols = 6, 9
	case "difficile":
		rows, cols = 7, 8
	default:
		http.Error(w, "Difficulté invalide", http.StatusBadRequest)
		return
	}

	grid := make([][]string, rows)
	for i := range grid {
		grid[i] = make([]string, cols)
	}

	models.CurrentGame = &models.Game{
		Player1:      player1,
		Player2:      player2,
		Grid:         grid,
		Turn:         0,
		Winner:       "",
		Difficulty:   difficulty,
		Player1Score: 0,
		Player2Score: 0,
	}

	http.Redirect(w, r, "/play", http.StatusSeeOther)
}

func PlayMove(w http.ResponseWriter, r *http.Request) {
	if models.CurrentGame == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost && models.CurrentGame.Winner == "" {
		columnStr := r.FormValue("column")
		if columnStr != "" {
			column, err := strconv.Atoi(columnStr)
			if err != nil || column < 0 || column >= len(models.CurrentGame.Grid[0]) {
				http.Error(w, "Colonne invalide", http.StatusBadRequest)
				return
			}
			symbol := ""
			if models.CurrentGame.Turn == 0 {
				symbol = "X"
			} else {
				symbol = "O"
			}
			for i := len(models.CurrentGame.Grid) - 1; i >= 0; i-- {
				if models.CurrentGame.Grid[i][column] == "" {
					models.CurrentGame.Grid[i][column] = symbol
					if CheckWin(models.CurrentGame.Grid, symbol) {
						if symbol == "X" {
							models.CurrentGame.Winner = models.CurrentGame.Player1
							models.CurrentGame.Player1Score++
						} else {
							models.CurrentGame.Winner = models.CurrentGame.Player2
							models.CurrentGame.Player2Score++
						}
					} else if IsDraw(models.CurrentGame.Grid) {
						models.CurrentGame.Winner = "Egalité"
					}
					models.CurrentGame.Turn = 1 - models.CurrentGame.Turn
					break
				}
			}
		}
	}

	tmpl := template.Must(template.New("game.html").Funcs(template.FuncMap{
		"seq":       Seq,
		"sub":       func(a, b int) int { return a - b },
		"cellClass": CellClass,
	}).ParseFiles("templates/game.html"))

	tmpl.Execute(w, models.CurrentGame)
}

func Rematch(w http.ResponseWriter, r *http.Request) {
	if models.CurrentGame == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	rows, cols := 6, 7
	switch models.CurrentGame.Difficulty {
	case "facile":
		rows, cols = 6, 7
	case "moyen":
		rows, cols = 6, 9
	case "difficile":
		rows, cols = 7, 8
	}

	grid := make([][]string, rows)
	for i := range grid {
		grid[i] = make([]string, cols)
	}

	models.CurrentGame.Grid = grid
	models.CurrentGame.Turn = 0
	models.CurrentGame.Winner = ""

	http.Redirect(w, r, "/play", http.StatusSeeOther)
}

func QuitGame(w http.ResponseWriter, r *http.Request) {
	models.CurrentGame = nil
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
