#!/bin/bash

# Build script for GoFetch WebAssembly

echo "Building GoFetch for WebAssembly..."

# Set WASM environment variables
export GOOS=js
export GOARCH=wasm

# Build the WASM binary
go build -o dist/gofetch.wasm ./cmd/wasm

# Copy wasm_exec.js from Go installation
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" dist/

echo "✓ Build complete!"
echo "  - WASM binary: dist/gofetch.wasm"
echo "  - JS runtime: dist/wasm_exec.js"
