#!/bin/bash

set -e

echo "Building Twitter CLI..."

# Build for current platform
go build -o twt

echo "✅ Build complete: ./twt"

# Optionally build for multiple platforms
read -p "Build for all platforms? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]
then
    echo "Building for multiple platforms..."
    
    # macOS
    GOOS=darwin GOARCH=amd64 go build -o dist/twt-darwin-amd64
    GOOS=darwin GOARCH=arm64 go build -o dist/twt-darwin-arm64
    
    # Linux
    GOOS=linux GOARCH=amd64 go build -o dist/twt-linux-amd64
    GOOS=linux GOARCH=arm64 go build -o dist/twt-linux-arm64
    
    # Windows
    GOOS=windows GOARCH=amd64 go build -o dist/twt-windows-amd64.exe
    
    echo "✅ All builds complete in ./dist/"
fi