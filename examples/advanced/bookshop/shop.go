package bookshop

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"
)

// New returns a bookshop that sells the given catalogue and takes payments
// through payments. It leaves Admin.Export unimplemented, which rest.Handler
// serves as 501 Not Implemented. It is safe for concurrent use.
func New(catalogue []Book, payments Payments) API {
	var (
		mu    sync.Mutex
		books = slices.Clone(catalogue)
	)
	find := func(isbn ISBN) (*Book, error) {
		i := slices.IndexFunc(books, func(b Book) bool { return b.ISBN == isbn })
		if i < 0 {
			return nil, Errors.NotFound
		}
		return &books[i], nil
	}
	var shop API
	shop.Books.Search = func(_ context.Context, query Query) ([]Book, error) {
		mu.Lock()
		defer mu.Unlock()
		var found []Book
		for _, b := range books {
			if query.Author != "" && !strings.EqualFold(b.Author, query.Author) {
				continue
			}
			if query.MaxPrice != 0 && b.Price > query.MaxPrice {
				continue
			}
			found = append(found, b)
		}
		return found, nil
	}
	shop.Books.Get = func(_ context.Context, isbn ISBN) (Book, error) {
		mu.Lock()
		defer mu.Unlock()
		book, err := find(isbn)
		if err != nil {
			return Book{}, err
		}
		return *book, nil
	}
	shop.Books.Order = func(ctx context.Context, order Order) (Receipt, error) {
		mu.Lock()
		defer mu.Unlock()
		book, err := find(order.ISBN)
		if err != nil {
			return Receipt{}, err
		}
		if order.Copies < 1 || order.Copies > book.Stock {
			return Receipt{}, Errors.OutOfStock
		}
		total := book.Price * Cents(order.Copies)
		charge, err := payments.Charge(ctx, total)
		if err != nil {
			// The wrapped error keeps the provider's details for the server's
			// log; Tokens.Redact keeps them from the caller.
			return Receipt{}, Errors.PaymentFailed.As(fmt.Errorf("charge %d: %w", total, err))
		}
		book.Stock -= order.Copies
		return Receipt{ISBN: book.ISBN, Copies: order.Copies, Total: total, Charge: charge}, nil
	}
	shop.Admin.Restock = func(_ context.Context, isbn ISBN, copies int) error {
		mu.Lock()
		defer mu.Unlock()
		book, err := find(isbn)
		if err != nil {
			return err
		}
		book.Stock += copies
		return nil
	}
	return shop
}
