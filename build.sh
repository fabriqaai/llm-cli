#!/bin/bash
set -e

VERSION=${VERSION:-"dev"}
COMMIT=${COMMIT:-"$(git rev-parse --short HEAD)"}
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS="-X main.Version=$VERSION -X main.Commit=$COMMIT -X main.Date=$BUILD_TIME -s -w"

echo "Building llm-cli..."
echo "Version: $VERSION"
echo "Commit: $COMMIT"
go build -ldflags "$LDFLAGS" -o llm-cli .
echo "Done: ./llm-cli"
