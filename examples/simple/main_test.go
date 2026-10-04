package main

import (
	"context"
	"fmt"
	"os"

	"github.com/iqhive/iqlink/api/cmdl"
)

// The implementation is ordinary Go and can be called directly.
func Example_direct() {
	calc := New()
	sum, err := calc.Add(context.Background(), 2, 40)
	fmt.Println(sum, err)
	// Output: 42 <nil>
}

// cmdl.System runs the same implementation with explicit arguments, which is
// how main behaves for "simple hello gopher" and "simple add 2 40".
func Example_commandLine() {
	for _, args := range [][]string{
		{"simple"},
		{"simple", "hello", "gopher"},
		{"simple", "add", "2", "40"},
	} {
		if err := (cmdl.System{Args: args, Stdout: os.Stdout}).Run(New()); err != nil {
			fmt.Println("error:", err)
		}
	}
	// Output:
	// is a small calculator and greeter.
	// hello <name>
	// add <a> <b>
	// Hello, gopher!
	// 42
}
