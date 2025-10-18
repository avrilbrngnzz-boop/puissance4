package main

import (
	"fmt"
	"net/http"
	"power4/routes"
)

func main() {
	http.HandleFunc("/", routes.StartPage)
	http.HandleFunc("/start", routes.StartGame)
	http.HandleFunc("/play", routes.PlayMove)
	http.HandleFunc("/rematch", routes.Rematch)
	http.HandleFunc("/quit", routes.QuitGame)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Serveur lancé sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
