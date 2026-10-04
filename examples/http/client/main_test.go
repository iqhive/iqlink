package main

import (
	"context"
	"net/http/httptest"
	"os"

	"github.com/iqhive/iqlink/api"
	"github.com/iqhive/iqlink/api/rest"
	"github.com/iqhive/iqlink/examples/http/notes"
)

// The client talks real HTTP to the server's handler on a local test server.
func Example() {
	handler, err := rest.Handler(nil, notes.InMemory())
	if err != nil {
		panic(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()

	client := api.Import[notes.API](rest.API, server.URL, server.Client())
	if err := run(context.Background(), client, os.Stdout); err != nil {
		panic(err)
	}
	// Output:
	// created #1 "buy milk"
	// created #2 "water the plants"
	// 2 notes
	// deleted #1
	// get #1: not found
	// empty note: a note needs some text
}
