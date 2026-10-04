package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/iqhive/iqlink/api"
	"github.com/iqhive/iqlink/api/cmdl"
	"github.com/iqhive/iqlink/api/rest"
	"github.com/iqhive/iqlink/api/stub"
)

// Greeting is the result of the demonstration's Greet operation.
type Greeting struct {
	Message string `json:"message"`
}

func (g Greeting) String() string { return g.Message }

// API is the single specification every demonstration links against.
type API struct {
	api.Specification `api:"Hello" cmd:"hello"
		is an offline README demonstration.`

	Greet func(ctx context.Context, name string) (Greeting, error) `rest:"GET /greet/{name=%v}" cmdl:"greet %v"
		returns a greeting for name.`
	Add func(ctx context.Context, a, b int) (int, error) `rest:"GET /add?a=%v&b=%v" cmdl:"add %v %v"
		returns the sum of a and b.`
}

func implementation() API {
	return API{
		Greet: func(_ context.Context, name string) (Greeting, error) { return Greeting{"hello, " + name}, nil },
		Add:   func(_ context.Context, a, b int) (int, error) { return a + b, nil },
	}
}

// exchange records one HTTP request made by the imported client.
type exchange struct {
	request string // method and request URI
	status  int
	body    string // compacted JSON response body
}

// inProcess serves the client's requests with the REST handler directly, so
// no listener, port or network takes part in generation.
type inProcess struct {
	handler http.Handler
	log     *[]exchange
}

func (t inProcess) RoundTrip(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	t.handler.ServeHTTP(w, req)
	res := w.Result()
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var body bytes.Buffer
	if err := json.Compact(&body, raw); err != nil {
		return nil, fmt.Errorf("%s %s: %w: %s", req.Method, req.URL, err, raw)
	}
	*t.log = append(*t.log, exchange{req.Method + " " + req.URL.RequestURI(), res.StatusCode, body.String()})
	res.Body = io.NopCloser(bytes.NewReader(raw))
	return res, nil
}

type fixtureResults struct {
	greeting Greeting
	sum      int
	wire     []exchange
	cliGreet string
	cliAdd   string
	stubErr  string
}

const notImplemented = "not implemented"

// executeFixture links one implementation through the rest, cmdl and stub
// linkers. Nothing random, timed or remote enters generation; it fails loudly
// rather than render invented results when the project changes.
func executeFixture() fixtureResults {
	ctx := context.Background()
	impl := implementation()
	var result fixtureResults

	handler, err := rest.Handler(nil, impl)
	if err != nil {
		panic(err)
	}
	client := api.Import[API](rest.API, "http://hello.invalid", &http.Client{Transport: inProcess{handler, &result.wire}})
	if result.greeting, err = client.Greet(ctx, "gopher"); err != nil {
		panic(err)
	}
	if result.sum, err = client.Add(ctx, 2, 40); err != nil {
		panic(err)
	}

	command := func(args ...string) string {
		out, err := cmdl.System{Args: append([]string{"hello"}, args...)}.Output(impl)
		if err != nil {
			panic(fmt.Sprintf("hello %s: %v", strings.Join(args, " "), err))
		}
		return strings.TrimSuffix(string(out), "\n")
	}
	result.cliGreet = command("greet", "gopher")
	result.cliAdd = command("add", "2", "40")

	stubbed := api.Import[API](stub.API, stub.TODO, errors.New(notImplemented))
	if _, err := stubbed.Greet(ctx, "gopher"); err == nil {
		panic("stub: Greet returned no error")
	} else {
		result.stubErr = err.Error()
	}
	return result
}
