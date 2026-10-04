package bookshop

import (
	"context"
	"errors"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/iqhive/iqlink/api"
)

// Role is what a token allows its holder to do.
type Role int

const (
	Customer Role = iota + 1 // may use Books
	Staff                    // may also use Admin
)

// Tokens is an api.Auth for rest.Handler that maps bearer tokens to roles.
// Log receives the full text of every error that Redact hides; nil means
// log.Default().
type Tokens struct {
	Roles map[string]Role
	Log   *log.Logger
}

type roleKey struct{}

// RoleOf returns the role that Authenticate stored in ctx.
func RoleOf(ctx context.Context) Role {
	role, _ := ctx.Value(roleKey{}).(Role)
	return role
}

// Authenticate identifies the caller by its bearer token before the arguments
// are decoded, and records its role in the context the function receives.
func (t Tokens) Authenticate(ctx context.Context, r *http.Request, _ api.Function) (context.Context, error) {
	token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	role := t.Roles[token]
	if !ok || role == 0 {
		return ctx, Errors.Unauthenticated
	}
	return context.WithValue(ctx, roleKey{}, role), nil
}

// Authorize runs after the arguments are decoded. fn.Path is the function's
// namespace path, so ["Admin"] for Admin.Restock.
func (t Tokens) Authorize(ctx context.Context, _ *http.Request, fn api.Function, _ []reflect.Value) error {
	if len(fn.Path) > 0 && fn.Path[0] == "Admin" && RoleOf(ctx) != Staff {
		return api.ErrAccessDenied
	}
	return nil
}

// Redact sees every error before it is sent. Errors the caller can act on
// (4xx) pass through. Server-side failures (5xx) are logged in full, and the
// caller only learns the status, so a payment provider's message or an
// internal error never leaves the server.
func (t Tokens) Redact(_ context.Context, err error) error {
	status := http.StatusInternalServerError
	var withStatus interface{ StatusHTTP() int }
	if errors.As(err, &withStatus) {
		status = withStatus.StatusHTTP()
	}
	if status < 500 {
		return err
	}
	logger := t.Log
	if logger == nil {
		logger = log.Default()
	}
	logger.Printf("%d: %v", status, err)
	return redacted(status)
}

type redacted int

func (r redacted) Error() string   { return http.StatusText(int(r)) }
func (r redacted) StatusHTTP() int { return int(r) }
