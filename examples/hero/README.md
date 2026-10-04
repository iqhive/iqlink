# iqlink README hero

This example uses [iqhive/banner v1.0.1](https://github.com/iqhive/banner) for
wordmark glyphs and spacing, the entire seeded intro sweep, themes, cell/frame
model, SVG rendering and terminal playback. It is its own module, so banner
never becomes a dependency of iqlink itself; the replacement `../..` makes the
demonstrations always run against this checkout's source.

## Run and regenerate

From the repository root:

```sh
make hero-svg
make hero
make hero-check
# Shared CLI flags, including themes and looping:
env GOWORK=off go -C examples/hero run . -loop -theme light
```

`docs/hero.svg` is generated. Edit the source, then regenerate it. Redirected
playback prints the complete project card without escapes. Reduced motion and
unsupported animation show that same useful card. The window is 836 × 498
with an 80 × 20 terminal grid; captions must fit 80 columns.

## Content and evidence

Header: **One Go API specification, linked to many ways of calling it.**

Points: **REST client and server · command line · native ABI · stub · xray**,
the linkers listed in the project Readme. They reveal individually, 200ms apart.

One `API` struct in `fixture.go` declares `Greet` and `Add` with `rest` and
`cmdl` tags. The demonstrations link that one implementation three ways:

1. `rest.Handler` serves it and `api.Import(rest.API, …)` calls it. The client's
   `http.Client` sends each request straight to the handler, so the requests,
   status codes and JSON bodies shown are real but no listener or network is used.
2. `cmdl.System` runs it with the arguments `hello greet gopher` and
   `hello add 2 40`, capturing standard output.
3. `api.Import(stub.API, …)` returns the configured error.

## Maintenance map

| Change | Edit |
| --- | --- |
| Header, points, installation command | Constants and `points` in `frames.go` |
| Demonstrated specification or linkers | `fixture.go` |
| Scene layout and captions | `demonstrations` in `frames.go` |
| Reading time | `build`: quick motion retains its holds; extra time is shared across settled results |
| Expected project results and artifact drift | `hero_test.go` |
| Intro, fonts, layout, themes, SVG CSS or playback | The shared banner repository |

The loop is exactly 25 seconds: the opening, both demonstration scenes, the
project card and a one-second IQ Hive vector logo. The hero uses banner's
`iqcaps` theme (`theme` in `frames.go`): an uppercase `IQLINK` wordmark whose
`IQ` carries the IQ Hive logo's yellow/orange/red gradient, with `LINK` and the
demonstrations in the bluegreen palette. `-theme default`, `-theme bluegreen`
and `-theme light` preview the alternatives.

`hero_test.go` fails when iqlink's behavior changes what the scenes show, or
when `docs/hero.svg` no longer matches the source. Update the expectations only
after checking that the new output is correct, then run `make hero-svg`.

## Change the shared package locally

Temporarily add a replacement inside this example module:

```sh
env GOWORK=off go -C examples/hero mod edit -replace github.com/iqhive/banner=../../../banner
make hero-svg hero-check
# Remove the override, select the published version, tidy, regenerate.
env GOWORK=off go -C examples/hero mod edit -dropreplace github.com/iqhive/banner
env GOWORK=off go -C examples/hero mod tidy
make hero-svg hero-check
```

Do not commit a sibling banner override. After visual changes, inspect direct
SVG and a responsive `<img>` at full and narrow widths, both demonstrations,
the final cards, the loop transition and reduced motion.
