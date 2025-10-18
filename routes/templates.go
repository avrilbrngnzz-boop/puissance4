package routes

import (
	"fmt"
	"html/template"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl := template.New(name)
	tmpl, err := tmpl.ParseFiles("templates/" + name)
	if err != nil {
		fmt.Println("Erreur chargement template :", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		fmt.Println("Erreur exécution template :", err)
		http.Error(w, "Erreur serveur", http.StatusInternalServerError)
	}
}
