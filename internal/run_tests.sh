#!/bin/bash

# Test runner script for internal packages
echo "Running tests for internal packages..."

# Change to the internal directory
cd "$(dirname "$0")"

# Run config tests
echo "Running config tests..."
cd config
go test -v .
cd ..

echo "All internal package tests completed!" 
