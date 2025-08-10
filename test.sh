#!/bin/bash

echo "Running Sleep No More Tests..."
echo ""

echo "[1/3] Running unit tests..."
go test -v -short ./...
if [ $? -ne 0 ]; then
    echo "Unit tests failed!"
    exit 1
fi

echo ""
echo "[2/3] Running integration tests..."
go test -v -tags=integration ./...
if [ $? -ne 0 ]; then
    echo "Integration tests failed!"
    exit 1
fi

echo ""
echo "[3/3] Running benchmarks..."
go test -bench=. -benchmem ./...
if [ $? -ne 0 ]; then
    echo "Benchmarks failed!"
    exit 1
fi

echo ""
