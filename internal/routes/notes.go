package routes

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"md-notes/internal/models"
	"net/http"

	_ "modernc.org/sqlite"
)

type ViewContext struct {
	Notes  []models.Note
	Note   models.Note
	Error  error
	Offset int
}

func renderView(w http.ResponseWriter, templateName string, context ViewContext) error {
	tmpl, err := template.ParseFiles("./internal/templates/layout.html", fmt.Sprintf("./internal/templates/notes/%s", templateName))
	if err != nil {
		return err
	}
	err = tmpl.Execute(w, context)
	return err
}

func getNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := models.GetNotes()

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Something Went Wrong", http.StatusInternalServerError)
		return
	}

	log.Println(notes)
	err = renderView(w, "index.html", ViewContext{Notes: notes})
	if err != nil {
		log.Printf("Error", err.Error())
	}
}

func newNote(w http.ResponseWriter, r *http.Request) {
	err := renderView(w, "new.html", ViewContext{})
	if err != nil {
		log.Printf("Error", err.Error())
	}
}

func createNotes(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to save", http.StatusInternalServerError)
	}
	note := models.Note{Title: r.FormValue("title"), Body: sql.NullString{}}
	note, err = models.CreateNote(note)

	if err != nil {
		http.Error(w, "Failed to save", http.StatusUnprocessableEntity)
		log.Println(err.Error())
		return
	}

	err = renderView(w, "show.html", ViewContext{Note: note})
	if err != nil {
		http.Error(w, "Not Found", http.StatusInternalServerError)
	}
}

func showNote(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	note, err := models.GetNote(id)

	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
	}

	err = renderView(w, "show.html", ViewContext{Note: note})
	if err != nil {
		http.Error(w, "Not Found", http.StatusInternalServerError)
	}
}
