# Internal Packages

This directory contains the internal packages for the complaint escalator system.

## Structure

- `config/` - Configuration management package
  - `config.go` - Configuration structs and loading functions
  - `config_test.go` - Tests for configuration functionality
- `email/` - Email client package
  - `email.go` - Azure Communication Services email client
- `notification/` - Notification client package
  - `notification.go` - Notification sending functionality
- `ai/` - AI text generation package
  - `ai.go` - AI-powered text generation

## Testing

### Configuration Tests
```bash
# Run config tests
cd config
go test -v .

# Or use the test runner script
chmod +x ../run_tests.sh
../run_tests.sh
```

### Test Configuration
Tests use the test configuration file located at `../cmd/server/config-test.yaml` which contains:
- Test intervals and backoff times
- Test Azure Communication Services credentials
- Test email configurations
- Safe test endpoints and domains

## Package Dependencies

- `config` - No internal dependencies
- `email` - Depends on `config` for configuration
- `notification` - No internal dependencies
- `ai` - No internal dependencies

## Notes

- All packages are designed to be testable with mock configurations
- Tests use relative paths to access test configuration files
- Package boundaries are clearly defined to avoid circular dependencies
- Each package can be tested independently 