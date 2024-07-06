# Storage

## Overview

This directory contains all storage-related code for the AMD.
Within the `storage` package we declare all storage interfaces.
This is also where we define the seed data for Apps and Components.

The implementations of the storage interfaces are kept in `storage/postgres`.

## Migrations

New storage entities and changes to existing ones should be made
via SQL migrations. See the [README](https://github.com/brankas/rbac/blob/master/usermgm/migrations/README.md) for that directory.

## Testing

We provide a utility function within `storage/postgres` (`NewTestStorage`)
for use within tests that require real storage
(e.g. service-level success scenario tests).

It creates its own database to prevent collision with tests from packages
other than the one it is currently being used in, and greatly simplifies
cleanup.

We typically provide a singleton per package:

```go
var	_cacheStorage *postgres.Storage

func TestMain(m *testing.M) {
	var teardown func()
	dbstring := os.Getenv("DATABASE_CONNECTION")
	if dbstring != "" {
		_cacheStorage, teardown = postgres.NewTestStorage(dbstring, filepath.Join("..", "..", "migrations", "sql"))
	}

	exitCode := m.Run()

	if teardown != nil {
		teardown()
	}

	os.Exit(exitCode)
}

func newTestStorage(tb testing.TB) *postgres.Storage {
	if testing.Short() {
		tb.Skip("skipping tests that use postgres on -short")
	}
	if _cacheStorage == nil {
		tb.Skip("database unset, skipping test")
	}

	return _cacheStorage
}

func TestFoo(t *testing.T) {
	t.Parallel()

	store := newTestStorage(t)

	// test proper
}
```
