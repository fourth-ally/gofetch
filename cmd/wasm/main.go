//go:build js && wasm
// +build js,wasm

package main

import (
	"syscall/js"

	"github.com/nkechr1s/gofetch/wasm"
)

func main() {
	// Log to console for debugging
	js.Global().Get("console").Call("log", "GoFetch WASM: main() started")

	// Expose GoFetch functions to JavaScript
	wasm.ExposeFunctions()

	// Log successful initialization
	js.Global().Get("console").Call("log", "GoFetch WASM: Functions exposed successfully")

	// Keep the program running
	select {}
}
