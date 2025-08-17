#!/bin/bash

# Test runner script for the complaint escalator server
echo "Running tests for complaint escalator..."

# Run all tests
echo "Running tests..."
go test -v ./...

echo "All tests completed!" 
