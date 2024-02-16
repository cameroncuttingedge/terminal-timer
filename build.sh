#!/bin/bash

APP_NAME="terminal-timer"
VERSION="1.0.0"

# macOS Intel build
echo "Building for macOS (Intel)..."
GOOS=darwin GOARCH=amd64 go build -o $APP_NAME
zip "${APP_NAME}-${VERSION}-macos-intel.zip" $APP_NAME
rm $APP_NAME

# macOS ARM64 build
echo "Building for macOS (ARM64)..."
GOOS=darwin GOARCH=arm64 go build -o $APP_NAME
zip "${APP_NAME}-${VERSION}-macos-arm64.zip" $APP_NAME
rm $APP_NAME

# Windows build
echo "Building for Windows (AMD64)..."
GOOS=windows GOARCH=amd64 go build -o "${APP_NAME}.exe"
zip "${APP_NAME}-${VERSION}-windows-amd64.zip" "${APP_NAME}.exe"
rm "${APP_NAME}.exe"

# Linux build
echo "Building for Linux (AMD64)..."
GOOS=linux GOARCH=amd64 go build -o $APP_NAME
zip "${APP_NAME}-${VERSION}-linux-amd64.zip" $APP_NAME
rm $APP_NAME

echo "Build and packaging process completed."
