package models

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Note struct {
	ID    int64
	Title string
	Body  sql.NullString
}

func GetNotes() ([]Note, error) {
	notes := []Note{}

	dsn := "file:///Users/tmobaird/Projects/markdown-notes/notes.sqlite3"
	pool, err := sql.Open("sqlite", dsn)
	if err != nil {
		return notes, err
	}
	defer pool.Close()

	pool.SetConnMaxLifetime(0)
	pool.SetMaxIdleConns(3)
	pool.SetMaxOpenConns(3)

	rows, err := pool.Query("select * from notes;")

	if err != nil {
		return notes, err
	} else {
		for rows.Next() {
			n := Note{}
			err = rows.Scan(&n.ID, &n.Title, &n.Body)
			if err != nil {
				return []Note{}, err
			}
			notes = append(notes, n)
		}
	}
	return notes, nil
}

func CreateNote(note Note) (Note, error) {
	dsn := "file:///Users/tmobaird/Projects/markdown-notes/notes.sqlite3"
	pool, err := sql.Open("sqlite", dsn)
	if err != nil {
		return Note{}, err
	}
	defer pool.Close()

	pool.SetConnMaxLifetime(0)
	pool.SetMaxIdleConns(3)
	pool.SetMaxOpenConns(3)

	result, err := pool.Exec(fmt.Sprintf("INSERT INTO notes (title) VALUES (\"%s\");", note.Title))
	if err != nil {
		return Note{}, err
	}
	id, err := result.LastInsertId()
	note.ID = id

	if err != nil {
		return Note{}, err
	}

	return note, nil
}

func GetNote(id string) (Note, error) {
	note := Note{}

	dsn := "file:///Users/tmobaird/Projects/markdown-notes/notes.sqlite3"
	pool, err := sql.Open("sqlite", dsn)
	if err != nil {
		return note, err
	}
	defer pool.Close()

	pool.SetConnMaxLifetime(0)
	pool.SetMaxIdleConns(3)
	pool.SetMaxOpenConns(3)

	rows := pool.QueryRow(fmt.Sprintf("select * from notes where id = '%s' limit 1;", id))

	err = rows.Scan(&note.ID, &note.Title, &note.Body)
	if err != nil {
		return Note{}, err
	}
	return note, nil
}

func UpdateNote(note Note) (Note, error) {
	dsn := "file:///Users/tmobaird/Projects/markdown-notes/notes.sqlite3"
	pool, err := sql.Open("sqlite", dsn)
	if err != nil {
		return Note{}, err
	}
	defer pool.Close()

	pool.SetConnMaxLifetime(0)
	pool.SetMaxIdleConns(3)
	pool.SetMaxOpenConns(3)

	_, err = pool.Exec(fmt.Sprintf("UPDATE notes SET title = \"%s\", body = \"%s\" WHERE id = %d;", note.Title, note.Body.String, note.ID))
	if err != nil {
		return Note{}, err
	}

	if err != nil {
		return Note{}, err
	}

	return note, nil
}
