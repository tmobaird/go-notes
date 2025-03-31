package routes

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"md-notes/internal/models"
	"net/http"
	"strconv"
	"strings"

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

	err = renderView(w, "index.html", ViewContext{Notes: notes})
	if err != nil {
		log.Println("Error", err.Error())
	}
}

func newNote(w http.ResponseWriter, r *http.Request) {
	err := renderView(w, "new.html", ViewContext{})
	if err != nil {
		log.Println("Error", err.Error())
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

	http.Redirect(w, r, fmt.Sprintf("/notes/%d", note.ID), http.StatusSeeOther)
}

func showNote(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	note, err := models.GetNote(id)

	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
	}

	note.Body = sql.NullString{String: strings.ReplaceAll(note.Body.String, "\r\n", "<br>")}
	err = renderView(w, "show.html", ViewContext{Note: note})
	if err != nil {
		http.Error(w, "Not Found", http.StatusInternalServerError)
	}
}

func editNote(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	note, err := models.GetNote(id)

	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
	}

	err = renderView(w, "edit.html", ViewContext{Note: note})
	if err != nil {
		http.Error(w, "Not Found", http.StatusInternalServerError)
	}
}

func updateNote(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to save", http.StatusInternalServerError)
	}
	intid, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
	}
	note := models.Note{ID: int64(intid), Title: r.FormValue("title"), Body: sql.NullString{String: r.FormValue("body"), Valid: true}}
	note, err = models.UpdateNote(note)

	if err != nil {
		http.Error(w, "Unprocessible Entity", http.StatusUnprocessableEntity)
	}
	http.Redirect(w, r, fmt.Sprintf("/notes/%d", note.ID), http.StatusSeeOther)
}
