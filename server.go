package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

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

var currentGame *Game

func seq(start, end int) []int {
	s := make([]int, end-start+1)
	for i := range s {
		s[i] = start + i
	}
	return s
}

func cellClass(cell string) string {
	switch cell {
	case "X":
		return "X"
	case "O":
		return "O"
	default:
		return "empty"
	}
}

func checkWin(grid [][]string, symbol string) bool {
	rows := len(grid)
	cols := len(grid[0])
	for i := 0; i < rows; i++ {
		for j := 0; j <= cols-4; j++ {
			if grid[i][j] == symbol && grid[i][j+1] == symbol &&
				grid[i][j+2] == symbol && grid[i][j+3] == symbol {
				return true
			}
		}
	}
	for i := 0; i <= rows-4; i++ {
		for j := 0; j < cols; j++ {
			if grid[i][j] == symbol && grid[i+1][j] == symbol &&
				grid[i+2][j] == symbol && grid[i+3][j] == symbol {
				return true
			}
		}
	}
	for i := 0; i <= rows-4; i++ {
		for j := 0; j <= cols-4; j++ {
			if grid[i][j] == symbol && grid[i+1][j+1] == symbol &&
				grid[i+2][j+2] == symbol && grid[i+3][j+3] == symbol {
				return true
			}
		}
	}
	for i := 3; i < rows; i++ {
		for j := 0; j <= cols-4; j++ {
			if grid[i][j] == symbol && grid[i-1][j+1] == symbol &&
				grid[i-2][j+2] == symbol && grid[i-3][j+3] == symbol {
				return true
			}
		}
	}
	return false
}

func isDraw(grid [][]string) bool {
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == "" {
				return false
			}
		}
	}
	return true
}

func startPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.New("start.html")
	tmpl, err := tmpl.ParseFiles("templates/start.html")
	if err != nil {
		fmt.Println("Erreur chargement start.html :", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "start.html", nil)
	if err != nil {
		fmt.Println("Erreur exécution start.html :", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
	}
}

func startGame(w http.ResponseWriter, r *http.Request) {
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

	currentGame = &Game{
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

func playMove(w http.ResponseWriter, r *http.Request) {
	if currentGame == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost && currentGame.Winner == "" {
		columnStr := r.FormValue("column")
		if columnStr != "" {
			column, err := strconv.Atoi(columnStr)
			if err != nil || column < 0 || column >= len(currentGame.Grid[0]) {
				http.Error(w, "Colonne invalide", http.StatusBadRequest)
				return
			}

			symbol := ""
			if currentGame.Turn == 0 {
				symbol = "X"
			} else {
				symbol = "O"
			}

			for i := len(currentGame.Grid) - 1; i >= 0; i-- {
				if currentGame.Grid[i][column] == "" {
					currentGame.Grid[i][column] = symbol
					if checkWin(currentGame.Grid, symbol) {
						if symbol == "X" {
							currentGame.Winner = currentGame.Player1
							currentGame.Player1Score++
						} else {
							currentGame.Winner = currentGame.Player2
							currentGame.Player2Score++
						}
					} else if isDraw(currentGame.Grid) {
						currentGame.Winner = "Egalité"
					}
					currentGame.Turn = 1 - currentGame.Turn
					break
				}
			}
		}
	}

	tmpl := template.Must(template.New("game.html").Funcs(template.FuncMap{
		"seq":       seq,
		"sub":       func(a, b int) int { return a - b },
		"cellClass": cellClass,
	}).ParseFiles("templates/game.html"))

	tmpl.Execute(w, currentGame)
}

func rematch(w http.ResponseWriter, r *http.Request) {
	if currentGame == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	rows, cols := 6, 7
	switch currentGame.Difficulty {
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

	currentGame.Grid = grid
	currentGame.Turn = 0
	currentGame.Winner = ""

	http.Redirect(w, r, "/play", http.StatusSeeOther)
}

func quitGame(w http.ResponseWriter, r *http.Request) {
	currentGame = nil
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/", startPage)
	http.HandleFunc("/start", startGame)
	http.HandleFunc("/play", playMove)
	http.HandleFunc("/rematch", rematch)
	http.HandleFunc("/quit", quitGame)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Serveur lancé sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
