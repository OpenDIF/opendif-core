# Shared Test Utilities

This package provides shared test utilities for unit testing across the opendif-core repository.

## Database Mocking

### SetupMockDB

Use `SetupMockDB` for unit tests that need to mock database interactions:

```go
import "github.com/gov-dx-sandbox/exchange/shared/testutils"

func TestMyFunction(t *testing.T) {
    db, mock, cleanup := testutils.SetupMockDB(t)
    defer cleanup()
    
    // Configure mock expectations
    mock.ExpectQuery("SELECT ...").WillReturnRows(...)
    
    // Use db in your tests
    // ...
    
    // Verify all expectations were met
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

## Guidelines

### Unit Tests
- **MUST** use `SetupMockDB` from this package
- **MUST NOT** use real database connections
- **MUST** run without external dependencies (database, network, etc.)
- **MUST** use `go test -short` to skip integration-style tests

### Integration Tests
- **CAN** use real database connections (PostgreSQL via Docker Compose)
- **SHOULD** be located in `tests/integration/` directory
- **SHOULD** use `testutils.SetupPostgresTestDB` from `tests/integration/testutils`
- **CAN** require external services (databases, APIs, etc.)

### Test Naming
- Unit tests: `TestFunctionName`
- Integration tests: `TestFunctionName_Integration` or located in `tests/integration/`

## Examples

### Unit Test (with mocks)
```go
// exchange/consent-engine/v1/services/consent_service_test.go
func TestCreateConsent(t *testing.T) {
    db, mock, cleanup := testutils.SetupMockDB(t)
    defer cleanup()
    
    mock.ExpectQuery("INSERT INTO ...").WillReturnRows(...)
    
    service := NewConsentService(db, "http://portal")
    // ... test logic
}
```

### Integration Test (with real DB)
```go
// tests/integration/consent/consent_test.go
func TestConsent_CreateAndRetrieve(t *testing.T) {
    db := testutils.SetupPostgresTestDB(t)
    // ... test with real database
}
```

