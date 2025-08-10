#!/bin/bash

# Test runner script for the complaint escalator server
echo "Running tests for complaint escalator server..."

# Change to the cmd/server directory
cd "$(dirname "$0")"

# Run all tests
echo "Running escalator tests..."
go test -v ./escalator_test.go

echo "Running server tests..."
go test -v ./server_test.go

echo "All tests completed!" 