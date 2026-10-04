package main

import (
	"context"
	"errors"
	"fmt"
	"net/http/httptest"
	"os"

	"github.com/iqhive/iqlink/api"
	"github.com/iqhive/iqlink/api/stub"
	"github.com/iqhive/iqlink/examples/advanced/bookshop"
)

// The whole program, against a local test server.
func Example() {
	handler, err := newHandler(os.Stdout)
	if err != nil {
		panic(err)
	}
	server := httptest.NewServer(handler)
	defer server.Close()
	run(context.Background(), server.URL, os.Stdout)
	// Output:
	// customer:
	//   GET /books?author=Ursula+K.+Le+Guin&max_price=2000 -> 200
	//     A Wizard of Earthsea, $9.99, 978-0-000-00001-1
	//     The Dispossessed, $15.99, 978-0-000-00002-2
	//   POST /orders -> 200
	//     bought 2 for $19.98, charge ch_1
	//   POST /orders -> 409
	//     error 409: Conflict
	//     server log: 502: charge 7500: card declined (provider code 51)
	//   POST /orders -> 502
	//     error 502: Bad Gateway
	//   POST /orders -> 404
	//     error 404: not found
	//   POST /admin/stock/978-0-000-00002-2/5 -> 403
	//     error 403: Forbidden
	// staff:
	//   POST /admin/stock/978-0-000-00002-2/5 -> 204
	//   GET /books/978-0-000-00002-2 -> 200
	//     The Dispossessed: 6 in stock
	//     server log: 501: not implemented
	//   GET /admin/export -> 501
	//     error 501: Not Implemented
	// anonymous:
	//   GET /books/978-0-000-00001-1 -> 401
	//     error 401: Unauthorized
}

// Without HTTP in between, callers get the declared error values themselves.
// Here the Payments dependency is a stub that fails every charge, so the shop
// can be tested without a payment provider.
func Example_stub() {
	payments := api.Import[bookshop.Payments](stub.API, stub.Testing, errors.New("payments are switched off"))
	shop := bookshop.New(catalogue, payments)

	_, err := shop.Books.Order(context.Background(), bookshop.Order{ISBN: "978-0-000-00004-4", Copies: 1})
	var shopErr bookshop.Error
	if errors.As(err, &shopErr) {
		fmt.Printf("HTTP %d: %v\n", shopErr.StatusHTTP(), err)
	}
	_, err = shop.Books.Get(context.Background(), "978-0-000-00009-9")
	fmt.Println(errors.Is(err, bookshop.Errors.NotFound))
	// Output:
	// HTTP 502: charge 1250: payments are switched off
	// true
}
