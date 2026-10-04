//go:build wasm

package wasm_test

import (
	"fmt"
	"testing"

	"github.com/iqhive/iqlink/api/wasm"
	"github.com/iqhive/iqlink/api/wasm/internal/example"
	"github.com/iqhive/iqlink/ffi"
)

var FFI = wasm.Import[ffi.API]()
var Example = wasm.Import[example.API]()

func TestExample(m *testing.T) {
	Example.HelloWorld()
	fmt.Println(Example.HostArch())
}

//go:wasmimport example HostArch
func host_arch() ffi.String

func BenchmarkDynamicHostArch(m *testing.B) {
	for m.Loop() {
		Example.HostArch()
	}
}

func BenchmarkWasmImportHostArch(m *testing.B) {
	for m.Loop() {
		host_arch()
	}
}
