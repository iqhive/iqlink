// Package notes is the API specification shared by the notes server and
// client. Neither side writes HTTP code: the server serves an implementation
// of API with rest.Handler, and the client links an empty API to the server
// with api.Import.
package notes

import (
	"context"

	"github.com/iqhive/iqlink/api"
)

// API stores short text notes. The www tag is the default address clients
// connect to when they don't give one.
type API struct {
	api.Specification `api:"Notes" www:"http://localhost:8080"
		stores short text notes.`

	Create func(ctx context.Context, draft Draft) (Note, error) `rest:"POST /notes"
		saves a new note and returns it with its ID.`
	List func(ctx context.Context) ([]Note, error) `rest:"GET /notes"
		returns every note, oldest first.`
	Get func(ctx context.Context, id ID) (Note, error) `rest:"GET /notes/{id=%v}"
		returns the note with the given ID.`
	Delete func(ctx context.Context, id ID) error `rest:"DELETE /notes/{id=%v}"
		removes the note with the given ID.`
}

// ID identifies a note.
type ID int

// Draft is a note that hasn't been saved yet.
type Draft struct {
	Text string `json:"text"`
}

// Note is a saved note.
type Note struct {
	ID   ID     `json:"id"`
	Text string `json:"text"`
}
