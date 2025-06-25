#!/usr/bin/env bash
set -euo pipefail

# Ensure Go is installed
if ! command -v go >/dev/null 2>&1; then
  echo "Go is required but was not found. Please install Go 1.24.3 from https://go.dev/dl/" >&2
  exit 1
fi

echo "Fetching Go dependencies..."
go mod download

echo "Downloading sample media files..."
./get-media.sh
