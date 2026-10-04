// Command server serves the notes API over HTTP.
//
//	go run ./examples/http/server -addr localhost:8080
//
// Then, in another terminal, run the client or use curl:
//
//	go run ./examples/http/client
//	curl -H 'Content-Type: application/json' -d '{"text":"buy milk"}' localhost:8080/notes
//	curl localhost:8080/notes/1
package main

import (
	"flag"
	"log"

	"github.com/iqhive/iqlink/api/rest"
	"github.com/iqhive/iqlink/examples/http/notes"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "address to listen on")
	flag.Parse()
	log.Printf("serving the notes API on http://%s", *addr)
	// A nil api.Auth lets every caller use every operation; see
	// examples/advanced for an authenticated server.
	log.Fatal(rest.ListenAndServe(*addr, nil, notes.InMemory()))
}
