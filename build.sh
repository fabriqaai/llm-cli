#!/bin/bash
set -e

VERSION=${VERSION:-"dev"}
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS="-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME"

echo "Building llm-cli..."
go build -ldflags "$LDFLAGS" -o llm-cli .
echo "Done: ./llm-cli"
