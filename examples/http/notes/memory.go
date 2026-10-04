package notes

import (
	"context"
	"net/http"
	"slices"
	"sync"
)

// Errors returned by InMemory. Their StatusHTTP method tells rest.Handler
// which status to answer with.
var (
	ErrNotFound error = statusError{http.StatusNotFound, "note not found"}
	ErrNoText   error = statusError{http.StatusBadRequest, "a note needs some text"}
)

type statusError struct {
	status  int
	message string
}

func (e statusError) Error() string   { return e.message }
func (e statusError) StatusHTTP() int { return e.status }

// InMemory returns an implementation of API that keeps notes in memory. It is
// safe for concurrent use.
func InMemory() API {
	var (
		mu    sync.Mutex
		next  ID = 1
		notes []Note
	)
	find := func(id ID) int {
		return slices.IndexFunc(notes, func(n Note) bool { return n.ID == id })
	}
	return API{
		Create: func(_ context.Context, draft Draft) (Note, error) {
			if draft.Text == "" {
				return Note{}, ErrNoText
			}
			mu.Lock()
			defer mu.Unlock()
			note := Note{ID: next, Text: draft.Text}
			next++
			notes = append(notes, note)
			return note, nil
		},
		List: func(context.Context) ([]Note, error) {
			mu.Lock()
			defer mu.Unlock()
			return slices.Clone(notes), nil
		},
		Get: func(_ context.Context, id ID) (Note, error) {
			mu.Lock()
			defer mu.Unlock()
			if i := find(id); i >= 0 {
				return notes[i], nil
			}
			return Note{}, ErrNotFound
		},
		Delete: func(_ context.Context, id ID) error {
			mu.Lock()
			defer mu.Unlock()
			if i := find(id); i >= 0 {
				notes = slices.Delete(notes, i, i+1)
				return nil
			}
			return ErrNotFound
		},
	}
}
