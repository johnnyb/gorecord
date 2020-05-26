package migrator

import (
	"database/sql"
	// "github.com/johnnyb/gorecord/gorec"
)

func MigrateWithDirectory() {
	// Go through the directory in canonical order
	// Look for files named `.up.migration`
	// Run UpMigrateWithFile
}

func UpMigrateWithFile() {
	// Make sure the migration table exists, create it if it doesn't (schema_migrations with varchar `version` field)
	// Take filename and find migration identifier
	// Read file for commands
	// Execute commands (in transaction?)
	// Mark the migration as done
}

func DownMigrateWithFile() {

}

func PerformUpMigration(db *sql.DB, identifier string, migration string) {

}

func PerformDownMigration(db *sql.DB, identifier string, migration string) {

}
