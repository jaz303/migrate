package migrate

import (
	"database/sql"
	"fmt"
)

type Migration func(*sql.DB) error

const (
	SchemaVersionTable = "schema_version"
)

func MustExec(db *sql.DB, sql string, args ...any) sql.Result {
	res, err := db.Exec(sql, args...)
	if err != nil {
		panic(err)
	}
	return res
}

func Migrate(db *sql.DB, migrations []Migration) error {
	initialVersion := getSchemaVersion(db)

	if initialVersion < 0 {
		if err := createSchemaVersionTable(db); err != nil {
			return err
		}
		initialVersion = 0
	}

	for idx, migration := range migrations {
		targetVersion := idx + 1
		if targetVersion > initialVersion {
			if err := runMigration(db, migration); err != nil {
				return fmt.Errorf("failed to migrate to version %d (%s)", targetVersion, err)
			} else if err := setSchemaVersion(db, targetVersion); err != nil {
				return err
			}
		}
	}

	return nil
}

func runMigration(db *sql.DB, m Migration) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case error:
				err = x
			default:
				err = fmt.Errorf("panic from migration callback: %v", x)
			}
		}
	}()
	err = m(db)
	return err
}

func getSchemaVersion(db *sql.DB) int {
	var version int
	if db.QueryRow("SELECT version FROM "+SchemaVersionTable).Scan(&version) != nil {
		return -1
	} else {
		return version
	}
}

func setSchemaVersion(db *sql.DB, version int) error {
	if _, err := db.Exec("UPDATE "+SchemaVersionTable+" SET version = ?", version); err != nil {
		return fmt.Errorf("failed to set schema version to %d (%v)", version, err)
	}
	return nil
}

func createSchemaVersionTable(db *sql.DB) error {
	if _, err := db.Exec("CREATE TABLE " + SchemaVersionTable + " (version INTEGER NOT NULL PRIMARY KEY)"); err != nil {
		return err
	}

	if _, err := db.Exec("INSERT INTO " + SchemaVersionTable + " (version) VALUES (0)"); err != nil {
		return err
	}

	return nil
}
