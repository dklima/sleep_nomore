#!/bin/bash

echo "Running Sleep No More Tests..."
echo ""

echo "[1/3] Running unit tests..."
echo "----------------------------------------"
go test -v -short ./...
if [ $? -ne 0 ]; then
    echo ""
    echo "Unit tests failed!"
    echo ""
    exit 1
fi
echo "Unit tests passed!"

echo ""
echo "[2/3] Running integration tests..."
echo "----------------------------------------"
go test -v -tags=integration ./...
if [ $? -ne 0 ]; then
    echo ""
    echo "Integration tests failed!"
    echo ""
    exit 1
fi
echo "Integration tests passed!"

echo ""
echo "[3/3] Running benchmarks..."
echo "----------------------------------------"
go test -bench=. -benchmem ./...
if [ $? -ne 0 ]; then
    echo ""
    echo "Benchmarks failed!"
    echo ""
    exit 1
fi
echo "Benchmarks completed successfully!"

echo ""
echo "Summary:"
echo "- Unit tests: PASS"
echo "- Integration tests: PASS"
echo "- Benchmarks: PASS"
echo "========================================"
