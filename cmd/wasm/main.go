//go:build js && wasm
// +build js,wasm

package main

import (
	"github.com/nkechr1s/gofetch/wasm"
)

func main() {
	// Expose GoFetch functions to JavaScript
	wasm.ExposeFunctions()

	// Keep the program running
	select {}
}
