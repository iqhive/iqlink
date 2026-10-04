// Package bookshop specifies a small online bookshop. Its API is organised
// into namespaces, declares the errors it can return along with their HTTP
// statuses, and depends on a second iqlink API, Payments.
package bookshop

import (
	"context"

	"github.com/iqhive/iqlink/api"
	"github.com/iqhive/iqlink/xyz"
)

// API is the bookshop. Nested structs are namespaces: Books is for anyone
// with a token and Admin for staff, which Tokens.Authorize enforces using the
// namespace path of each function.
type API struct {
	api.Specification `api:"Bookshop" www:"http://localhost:8080"
		sells books.`

	// Register documents every Error case as a possible response, with its
	// HTTP status, in the API's structure and generated documentation.
	_ api.Register[error, Error]

	Books struct {
		Search func(ctx context.Context, query Query) ([]Book, error) `rest:"GET /books?%v"
			lists the books that match every field set in query.`
		Get func(ctx context.Context, isbn ISBN) (Book, error) `rest:"GET /books/{isbn=%v}"
			returns one book.`
		Order func(ctx context.Context, order Order) (Receipt, error) `rest:"POST /orders"
			charges for and reserves copies of a book.`
	}
	Admin struct {
		Restock func(ctx context.Context, isbn ISBN, copies int) error `rest:"POST /admin/stock/{isbn=%v}/{copies=%v}"
			adds copies of a book to the stock.`
		Export func(ctx context.Context) ([]Book, error) `rest:"GET /admin/export"
			returns the whole catalogue; not implemented yet.`
	}
}

// Error enumerates the bookshop's errors. The http tag is the status that
// rest.Handler answers with, and the lines after it document the case.
type Error api.Error[struct {
	Unauthenticated Error `http:"401"
		the request has no valid bearer token.`
	NotFound Error `http:"404"
		there is no book with that ISBN.`
	OutOfStock Error `http:"409"
		there are not enough copies in stock.`
	PaymentFailed xyz.Case[Error, error] `http:"502"
		the payment provider could not take the payment; it wraps the
		provider's error.`
}]

// Errors holds each Error case, as in Errors.NotFound.
var Errors = xyz.AccessorFor(Error.Values)

// ISBN identifies a book.
type ISBN string

// Cents is an amount of money in hundredths of a dollar.
type Cents int

// Query filters Books.Search. Its fields become query parameters, and zero
// fields are left out.
type Query struct {
	Author   string `json:"author,omitempty"`
	MaxPrice Cents  `json:"max_price,omitempty"`
}

type Book struct {
	ISBN   ISBN   `json:"isbn"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Price  Cents  `json:"price"`
	Stock  int    `json:"stock"`
}

type Order struct {
	ISBN   ISBN `json:"isbn"`
	Copies int  `json:"copies"`
}

type Receipt struct {
	ISBN   ISBN     `json:"isbn"`
	Copies int      `json:"copies"`
	Total  Cents    `json:"total"`
	Charge ChargeID `json:"charge"`
}

// Payments is the payment provider the bookshop depends on. It is another
// iqlink API, so the shop can be given a REST client for a real provider, a
// Go implementation, or a stub, without changing.
type Payments struct {
	api.Specification `api:"Payments"`

	Charge func(ctx context.Context, amount Cents) (ChargeID, error) `rest:"POST /charges/{amount=%v}"
		takes a payment and returns its ID.`
}

// ChargeID identifies a payment.
type ChargeID string
