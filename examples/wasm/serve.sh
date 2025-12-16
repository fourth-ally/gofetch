#!/bin/bash

# Serve the WASM demo locally

echo "Starting local server for WASM demo..."
echo "Navigate to: http://localhost:8080"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

# Check if Python 3 is available
if command -v python3 &> /dev/null; then
    python3 -m http.server 8080
elif command -v python &> /dev/null; then
    python -m SimpleHTTPServer 8080
else
    echo "Error: Python is not installed"
    exit 1
fi
