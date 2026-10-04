// Command client calls the notes API served by examples/http/server.
//
//	go run ./examples/http/client                          # the www tag's address
//	go run ./examples/http/client -url http://host:8080    # another server
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/iqhive/iqlink/api"
	"github.com/iqhive/iqlink/api/rest"
	"github.com/iqhive/iqlink/examples/http/notes"
)

func main() {
	url := flag.String("url", "", "server address (default: the API's www tag)")
	flag.Parse()
	// api.Import fills in every function of notes.API with one that sends the
	// HTTP request its rest tag describes and decodes the response.
	client := api.Import[notes.API](rest.API, *url, http.DefaultClient)
	if err := run(context.Background(), client, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

// run uses the API like any Go value; it can't tell that each call is an HTTP
// request.
func run(ctx context.Context, client notes.API, w io.Writer) error {
	for _, text := range []string{"buy milk", "water the plants"} {
		note, err := client.Create(ctx, notes.Draft{Text: text})
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "created #%d %q\n", note.ID, note.Text)
	}
	all, err := client.List(ctx)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%d notes\n", len(all))

	first := all[0].ID
	if err := client.Delete(ctx, first); err != nil {
		return err
	}
	fmt.Fprintf(w, "deleted #%d\n", first)
	if _, err := client.Get(ctx, first); err != nil {
		fmt.Fprintf(w, "get #%d: %v\n", first, err)
	}
	if _, err := client.Create(ctx, notes.Draft{}); err != nil {
		fmt.Fprintf(w, "empty note: %v\n", err)
	}
	return nil
}
