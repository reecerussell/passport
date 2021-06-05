#!/bin/bash

set -e

if [[ "$CONFIG_DIR" != '' ]]; then
    echo "Config Directory: $CONFIG_DIR"
fi

if [[ "$VERSION" != '' ]]; then
    echo "Version: $VERSION"
fi

export GOARCH=amd64
export CGO_ENABLED=0

echo "Building..."
go build -o passport.exe \
    -ldflags "-X main.configDir=$CONFIG_DIR -X main.version=$VERSION" \
    cmd/main.go

echo "Done."