package migrator

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/johnnyb/gorecord/gorec"
	"sort"
)

// MigrationFunction is the type of function that is used for migrations
type MigrationFunction func(conn gorec.Querier) error

// Migration holds the up and down migration functions
type Migration struct {
	Version      string
	UpMigrator   MigrationFunction
	DownMigrator MigrationFunction
}

var registeredMigrations []Migration = []Migration{}

// RegisterMigration registers a migration to be run by the migrator
func RegisterMigration(v string, up MigrationFunction, down MigrationFunction) {
	registeredMigrations = append(registeredMigrations, Migration{
		Version:      v,
		UpMigrator:   up,
		DownMigrator: down,
	})
}

// DownMigrationNotPermitted is the standard down function for not allowing a down migration
func DownMigrationNotPermitted(conn gorec.Querier) error {
	return errors.New("Down migration not permitted")
}

// IsMigrationCompleted checks to see if a given migration has occurred
func IsMigrationCompleted(db gorec.Querier, version string) bool {
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return true // If we don't know, pretend it has been done
	}
	defer rows.Close()

	return rows.Next()
}

// Sorts the migrations by version name
func prepareMigrations() {
	sort.Slice(registeredMigrations, func(i, j int) bool {
		return registeredMigrations[i].Version < registeredMigrations[j].Version
	})

	migrationList := map[string]bool{}

	// Validate that we don't have any duplicates
	for _, migration := range registeredMigrations {
		if migrationList[migration.Version] {
			panic("Migration " + migration.Version + " registered twice")
		}
	}
}

// MigrateRegisteredMigrations is the automagic function to do all the necessary migrations
func MigrateRegisteredMigrations() error {
	prepareMigrations()
	return PerformUpMigrationsToVersion(gorec.GlobalConnection, registeredMigrations, registeredMigrations[len(registeredMigrations)-1].Version, true)
}

// PerformUpMigrationsToVersion uses the given connection to perform all of the migrations until it gets to the version specified (including the version specified)
func PerformUpMigrationsToVersion(conn *sql.DB, migrations []Migration, targetVersion string, useTransactions bool) error {
	AutoCreateSchemaTableIfNecessary(conn) // ignore errors because if the table already exists we don't care

	for _, migration := range registeredMigrations {
		if IsMigrationCompleted(conn, migration.Version) {
			if useTransactions {
				err := gorec.WithDBTransaction(conn, func(tx *sql.Tx) error {
					return UpMigrate(tx, migration)
				})
				if err != nil {
					return err
				}
			} else {
				err := UpMigrate(conn, migration)
				if err != nil {
					return err
				}
			}
		}

		if migration.Version == targetVersion {
			break
		}
	}

	return nil
}

// AutoCreateSchemaTableIfNecessary generates the schema table if it doesn't already exist.  Note that this will generate an error if the table already exists, so don't run this inside a transaction.
func AutoCreateSchemaTableIfNecessary(conn gorec.Querier) {
	conn.Exec("CREATE TABLE schema_migrations (version varchar(255) PRIMARY KEY)")
}

// UpMigrate runs a given migration up
func UpMigrate(conn gorec.Querier, migration Migration) error {
	fmt.Printf("Migrating Up: %s", migration.Version)
	err := migration.UpMigrator(conn)
	if err != nil {
		return err
	}
	_, err = conn.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", migration.Version)
	return err
}

// DownMigrate runs a given migration up
func DownMigrate(conn gorec.Querier, migration Migration) error {
	err := migration.DownMigrator(conn)
	if err != nil {
		return err
	}
	_, err = conn.Exec("DELETE FROM schema_migrations WHERE version = $1", migration.Version)
	return err
}
