#!/usr/bin/env bash
set -e

BIN_DIR="./bin"
CMD_DIR="./cmd/vunat"

mkdir -p "$BIN_DIR"

echo "Building vunat-cli for Linux..."
GOOS=linux GOARCH=amd64 go build -o "$BIN_DIR/vunat-linux-amd64" "$CMD_DIR"

echo "Building vunat-cli for Windows..."
GOOS=windows GOARCH=amd64 go build -o "$BIN_DIR/vunat-windows-amd64.exe" "$CMD_DIR"

echo "Build successful:"
echo "  - Linux: $BIN_DIR/vunat-linux-amd64"
echo "  - Windows: $BIN_DIR/vunat-windows-amd64.exe"
