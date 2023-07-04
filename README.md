# migrate

A simple Go library for handling forward-only database migrations, intended for use in embedded scenarios (e.g. sqlite).

## How to use

First, define your migrations as an array of callbacks:

```go
var migrations = []migrate.Migration{
    func(db *sql.DB) error {
        // MustExec panics on error; the migrator will recover
        // from these and return an error.
        migrate.MustExec(`
        
        `)
    },
    func(db *sql.DB) error {
        // You can also return an error instead of panicking...
        _, err := db.Exec("...")
        return err
    },
}
```

Then run your migrations:

```go
migrate.Migrate(db, migrations)
```