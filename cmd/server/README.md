# Server Command

This directory contains the server command and tests for the complaint escalator system.

## Structure

- `main.go` - Main server entry point
- `config-test.yaml` - Test configuration file with test values for all services
- `email_test.go` - Tests for the email package
- `escalator_test.go` - Tests for the main escalator functionality
- `run_tests.sh` - Script to run all tests

## Test Configuration

The `config-test.yaml` file contains test values for:
- **Interval**: 5 minutes (for faster testing)
- **Backoff**: 2 minutes (for faster testing)
- **Template**: Test complaint template
- **Channels**: email and notification
- **ACS**: Test Azure Communication Services configuration
- **Email**: Test email configuration with CC, BCC, and Reply-To

## Running Tests

### Individual Tests
```bash
# Test email functionality
go test -v ./email_test.go

# Test escalator functionality
go test -v ./escalator_test.go
```

### All Tests
```bash
# Make the script executable
chmod +x run_tests.sh

# Run all tests
./run_tests.sh
```

**Note**: Configuration tests are located in `./internal/config/` directory where they belong.

## Test Values

All test values are configured in `config-test.yaml` and are designed to:
- Use test domains and endpoints
- Have realistic but safe test credentials
- Include all optional fields for comprehensive testing
- Use shorter intervals for faster test execution

## Notes

- Tests use the test configuration instead of hardcoded values
- Email tests will fail on actual sending due to test credentials (this is expected)
- The test configuration is separate from production configuration
- All tests validate the configuration loading and structure, not external service calls
- Tests are in the same package as main.go for easier access to internal functions
- Configuration tests are properly located in the internal/config package 