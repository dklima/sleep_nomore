#!/bin/bash

echo "Downloading dependencies..."
go mod download

echo ""
echo "Building for Windows x86_64..."
GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui -s -w" -o sleepnomore_amd64.exe main.go

echo ""
echo "Building for Windows ARM64..."
GOOS=windows GOARCH=arm64 go build -ldflags="-H windowsgui -s -w" -o sleepnomore_arm64.exe main.go

echo ""
echo "Build complete!"
echo "Generated binaries:"
echo "- sleepnomore_amd64.exe (x86_64)"
echo "- sleepnomore_arm64.exe (ARM64)"