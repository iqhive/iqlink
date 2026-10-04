// Command simple is the smallest useful iqlink program: an API specification,
// a Go implementation of it, and the cmdl linker turning it into a command
// line.
//
//	go run ./examples/simple                 # prints the usage
//	go run ./examples/simple hello gopher    # Hello, gopher!
//	go run ./examples/simple add 2 40        # 42
package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"github.com/iqhive/iqlink/api"
	"github.com/iqhive/iqlink/api/cmdl"
)

// API is the specification. Each function field is an operation; its cmdl tag
// says how the command line addresses it, and the lines after the tag
// document it.
type API struct {
	api.Specification `cmd:"simple"
		is a small calculator and greeter.`

	Hello func(ctx context.Context, name string) (string, error) `cmdl:"hello %v"
		greets name.`
	Add func(ctx context.Context, a, b int) (int, error) `cmdl:"add %v %v"
		adds two whole numbers.`
}

//go:embed main.go
var source embed.FS

// Source lets linkers read the parameter names above, so the usage text says
// "hello <name>" rather than "hello <string>".
func (API) Source() fs.FS { return source }

// New returns the implementation. It is plain Go: nothing in it knows that it
// will be called from a command line.
func New() API {
	return API{
		Hello: func(_ context.Context, name string) (string, error) {
			return fmt.Sprintf("Hello, %s!", name), nil
		},
		Add: func(_ context.Context, a, b int) (int, error) {
			return a + b, nil
		},
	}
}

func main() { cmdl.Main(New()) }
