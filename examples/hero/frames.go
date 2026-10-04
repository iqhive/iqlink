package main

import (
	"fmt"
	"time"

	"github.com/iqhive/banner"
)

// All styling, wordmark layout, opening animation, SVG and playback live in banner.
// This file owns project text, deterministic demonstrations and reading holds.
type frame = banner.Frame

const (
	plain       = banner.Plain
	dim         = banner.Dim
	muted       = banner.Muted
	bright      = banner.Bright
	lit         = banner.Literal
	hot         = banner.Changed
	on          = banner.On
	prompt      = banner.Prompt
	projectName = "iqlink"
	tagline     = "One Go API specification, linked to many ways of calling it."
	subline     = "REST client and server · command line · native ABI · stub · xray"
	install     = "go get github.com/iqhive/iqlink@latest"
)

// theme gives the wordmark uppercase glyphs with the IQ Hive logo's gradient
// on its IQ prefix.
var theme = banner.IQCapsTheme()

var points = []string{"REST client and server", "command line", "native ABI", "stub", "xray"}

type film struct {
	cur, card frame
	frames    []frame
}

func (m *film) cut(ms int) {
	f := m.cur
	f.Hold = time.Duration(ms) * time.Millisecond
	m.frames = append(m.frames, f)
}

// build executes the project fixtures. The seed affects decorative digits only.
func build(seed uint64) []frame {
	m := &film{}
	opening, card, err := banner.Opening(banner.Card{
		Word: projectName, Tagline: tagline, Points: points, Install: install,
		Wordmark: banner.WordmarkOptions{Seed: seed}, ReadHold: 2500 * time.Millisecond,
		TaglineReveal: banner.TaglineWithReveal, PointInterval: 200 * time.Millisecond,
		Theme: theme,
	})
	if err != nil {
		panic(err)
	}
	m.frames, m.card = opening, card
	m.demonstrations()
	m.cur = m.card
	m.cut(2500)
	m.cur = frame{HiveLogo: true}
	m.cur.Center(14, "IQ Hive", bright)
	m.cur.Center(16, "iqhive.com", muted)
	m.cut(1000)
	// Keep all scenes and quick motion; distribute extra reading time among
	// settled demonstration results and the final project card.
	extra := 25*time.Second - banner.Duration(m.frames)
	if extra < 0 {
		panic("hero: scene timing exceeds 25 seconds")
	}
	var holds []int
	for i := len(opening); i < len(m.frames)-1; i++ {
		if m.frames[i].Hold >= time.Second {
			holds = append(holds, i)
		}
	}
	for n, i := range holds {
		share := extra / time.Duration(len(holds)-n)
		m.frames[i].Hold += share
		extra -= share
	}
	return m.frames
}

// animation keeps generator failures visible instead of writing guessed output.
func animation(seed uint64) (a banner.Animation, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("hero fixture: %v", p)
		}
	}()
	frames := build(seed)
	a = banner.Animation{Frames: frames, StaticFrame: len(frames) - 2,
		Title:       projectName + " — README demonstration",
		Description: "One Go API specification served and called over REST, run as a command line and stubbed, using offline fixtures.",
		Command:     "make hero",
		Theme:       theme,
	}
	return a, a.Validate()
}

// typed types input at (x, y), keeping a shell prompt in front of it when
// shell is set, and holds the completed line for ms.
func (m *film) typed(x, y int, input string, shell bool, ms int) {
	if shell {
		m.cur.Text(x, y, "$", prompt)
		x += 2
	}
	steps, done, err := banner.TypeLine(m.cur, x, y, input, banner.TypingOptions{TextStyle: lit, StepHold: 35 * time.Millisecond})
	if err != nil {
		panic(err)
	}
	for _, f := range append(steps, done) {
		if shell {
			f.Text(x-2, y, "$", prompt)
		}
		m.cur = f
		if f.Hold > 0 {
			m.frames = append(m.frames, f)
		}
	}
	m.cut(ms)
}

func (m *film) demonstrations() {
	result := executeFixture()
	wire := func(y int, e exchange) {
		m.cur.Text(6, y, fmt.Sprintf("%-20s %d  %s", e.request, e.status, e.body), dim)
	}

	m.cur = frame{}
	m.cur.Text(3, 1, "01 / Go specification -> rest.Handler -> typed client", muted)
	m.cut(150)
	m.cur.Text(4, 3, "Greet func(ctx context.Context, name string) (Greeting, error)", plain)
	m.cur.Text(10, 4, `rest:"GET /greet/{name=%v}"   cmdl:"greet %v"`, dim)
	m.cut(150)
	m.cur.Text(4, 6, "Go specification", bright)
	m.cur.Text(21, 6, "────> rest.Handler ────> api.Import(rest.API)", on)
	m.cut(450)
	m.typed(4, 9, `client.Greet(ctx, "gopher")`, false, 200)
	wire(10, result.wire[0])
	m.cur.Text(6, 11, fmt.Sprintf("Greeting{Message: %q}", result.greeting.Message), bright)
	m.cut(1200)
	m.typed(4, 13, "client.Add(ctx, 2, 40)", false, 200)
	wire(14, result.wire[1])
	m.cur.Text(6, 15, fmt.Sprint(result.sum), on)
	m.cur.Text(4, 18, "In-process HTTP through the real REST linker / offline fixture", muted)
	m.cut(2600)

	m.cur = frame{}
	m.cur.Text(3, 1, "02 / The same struct as a command line, then as a stub", muted)
	m.cut(150)
	m.typed(4, 4, "hello greet gopher", true, 200)
	m.cur.Text(4, 5, result.cliGreet, bright)
	m.cut(1000)
	m.typed(4, 7, "hello add 2 40", true, 200)
	m.cur.Text(4, 8, result.cliAdd, on)
	m.cut(1000)
	m.typed(4, 11, "hello := api.Import[API](stub.API, stub.TODO, err)", false, 200)
	m.typed(4, 12, `hello.Greet(ctx, "gopher")`, false, 200)
	m.cur.Text(6, 13, "error: "+result.stubErr, hot)
	m.cur.Text(4, 15, "Go specification", bright)
	m.cur.Text(21, 15, "────> rest  cmdl  link  stub  xray", on)
	m.cur.Text(4, 18, "cmdl.System parses real arguments / stub returns zero values and err", muted)
	m.cut(2600)
}
