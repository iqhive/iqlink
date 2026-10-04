// Command advanced serves the bookshop API on a local port and calls it as
// three different clients, printing each HTTP exchange. It shows:
//
//   - namespaces (Books, Admin) and a query-string struct (Books.Search);
//
//   - declared errors with HTTP statuses (bookshop.Error);
//
//   - bearer-token authentication, per-namespace authorization and redaction
//     of server-side errors with an api.Auth (bookshop.Tokens);
//
//   - an API that depends on another iqlink API (bookshop.Payments);
//
//   - an operation without an implementation, served as 501;
//
//   - client-side request interception with rest.Interceptor and per-client
//     headers with rest.Header.
//
//     go run ./examples/advanced
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/iqhive/iqlink/api"
	"github.com/iqhive/iqlink/api/rest"
	"github.com/iqhive/iqlink/examples/advanced/bookshop"
)

// The ISBNs are placeholders.
var catalogue = []bookshop.Book{
	{ISBN: "978-0-000-00001-1", Title: "A Wizard of Earthsea", Author: "Ursula K. Le Guin", Price: 999, Stock: 3},
	{ISBN: "978-0-000-00002-2", Title: "The Dispossessed", Author: "Ursula K. Le Guin", Price: 1599, Stock: 1},
	{ISBN: "978-0-000-00003-3", Title: "Collected Hainish Novels", Author: "Ursula K. Le Guin", Price: 7500, Stock: 2},
	{ISBN: "978-0-000-00004-4", Title: "Kindred", Author: "Octavia E. Butler", Price: 1250, Stock: 4},
}

// payments stands in for a payment provider: it declines anything over $50.
func payments() bookshop.Payments {
	charges := 0
	return bookshop.Payments{
		Charge: func(_ context.Context, amount bookshop.Cents) (bookshop.ChargeID, error) {
			if amount > 5000 {
				return "", errors.New("card declined (provider code 51)")
			}
			charges++
			return bookshop.ChargeID(fmt.Sprintf("ch_%d", charges)), nil
		},
	}
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	handler, err := newHandler(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	server := &http.Server{Handler: handler}
	go server.Serve(listener)
	defer server.Shutdown(context.Background())

	run(context.Background(), "http://"+listener.Addr().String(), os.Stdout)
}

// newHandler serves the bookshop with token authentication. Errors hidden
// from callers are logged to w.
func newHandler(w io.Writer) (http.Handler, error) {
	tokens := bookshop.Tokens{
		Roles: map[string]bookshop.Role{"alice-token": bookshop.Customer, "bob-token": bookshop.Staff},
		Log:   log.New(w, "    server log: ", 0),
	}
	return rest.Handler(tokens, bookshop.New(catalogue, payments()))
}

// run calls the bookshop served at url.
func run(ctx context.Context, url string, w io.Writer) {
	// Each client is a bookshop.API linked to the server. rest.Header gives a
	// client its own Authorization header.
	connect := func(token string) bookshop.API {
		client := http.DefaultClient
		if token != "" {
			client = rest.Header("Authorization", "Bearer "+token)
		}
		return api.Import[bookshop.API](rest.API, url, client)
	}
	// The interceptor sees every request the clients send with this context.
	ctx = rest.Interceptor(ctx, func(req *http.Request, send func(*http.Request) (*http.Response, error)) (*http.Response, error) {
		resp, err := send(req)
		if err == nil {
			fmt.Fprintf(w, "  %s %s -> %d\n", req.Method, req.URL.RequestURI(), resp.StatusCode)
		}
		return resp, err
	})
	failed := func(err error) {
		var withStatus interface{ StatusHTTP() int }
		if errors.As(err, &withStatus) {
			fmt.Fprintf(w, "    error %d: %v\n", withStatus.StatusHTTP(), err)
		} else {
			fmt.Fprintf(w, "    error: %v\n", err)
		}
	}

	fmt.Fprintln(w, "customer:")
	alice := connect("alice-token")
	books, err := alice.Books.Search(ctx, bookshop.Query{Author: "Ursula K. Le Guin", MaxPrice: 2000})
	if err != nil {
		failed(err)
	}
	for _, b := range books {
		fmt.Fprintf(w, "    %s, %s, %s\n", b.Title, dollars(b.Price), b.ISBN)
	}
	for _, order := range []bookshop.Order{
		{ISBN: "978-0-000-00001-1", Copies: 2}, // in stock and affordable
		{ISBN: "978-0-000-00002-2", Copies: 5}, // only one copy in stock
		{ISBN: "978-0-000-00003-3", Copies: 1}, // declined by the provider
		{ISBN: "978-0-000-00009-9", Copies: 1}, // no such book
	} {
		receipt, err := alice.Books.Order(ctx, order)
		if err != nil {
			failed(err)
			continue
		}
		fmt.Fprintf(w, "    bought %d for %s, charge %s\n", receipt.Copies, dollars(receipt.Total), receipt.Charge)
	}
	if err := alice.Admin.Restock(ctx, "978-0-000-00002-2", 5); err != nil {
		failed(err)
	}

	fmt.Fprintln(w, "staff:")
	bob := connect("bob-token")
	if err := bob.Admin.Restock(ctx, "978-0-000-00002-2", 5); err != nil {
		failed(err)
	}
	if book, err := bob.Books.Get(ctx, "978-0-000-00002-2"); err != nil {
		failed(err)
	} else {
		fmt.Fprintf(w, "    %s: %d in stock\n", book.Title, book.Stock)
	}
	if _, err := bob.Admin.Export(ctx); err != nil {
		failed(err)
	}

	fmt.Fprintln(w, "anonymous:")
	if _, err := connect("").Books.Get(ctx, "978-0-000-00001-1"); err != nil {
		failed(err)
	}
}

func dollars(c bookshop.Cents) string { return fmt.Sprintf("$%d.%02d", c/100, c%100) }
