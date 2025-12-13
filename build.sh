#!/usr/bin/env bash
set -e

BIN_DIR="./bin"
CMD_DIR="./cmd/vunat"

mkdir -p "$BIN_DIR"

echo "Building vunat-cli..."

go build -o "$BIN_DIR/vunat" "$CMD_DIR"

echo "Build successful: $BIN_DIR/vunat"
