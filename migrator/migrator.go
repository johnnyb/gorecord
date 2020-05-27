package migrator

import (
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"database/sql"
	"github.com/johnnyb/gorecord/gorec"
)

func IsMigrationCompleted(db *sql.DB, version string) bool {
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return true // If we don't know, pretend it has been done
	}
	defer rows.Close()

	return rows.Next()
}

func MigrateWithDirectory(db *sql.DB, dirname string, useTransactions bool)error {
	dir, err := os.Open(dirname)
	if err != nil {
		return err
	}

	AutoCreateSchemaTableIfNecessary(db) // ignore errors because if the table already exists we don't care

	names, err := dir.Readdirnames(0)
	sort.Strings(names)

	for _, name := range names {
		if strings.HasSuffix(name, ".up.sql") {
			version := name[:len(name) - 7]
				if IsMigrationCompleted(db, version) {
				if useTransactions {
					err = gorec.WithDBTransaction(db, func(tx *sql.Tx) error {
						return UpMigrateWithFile(version, dirname, name, tx)
					})
				} else {
					err = UpMigrateWithFile(version, dirname, name, db)
				}
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func AutoCreateSchemaTableIfNecessary(conn gorec.Querier) {
	conn.Exec("CREATE TABLE schema_migrations (version varchar(255) PRIMARY KEY)")
}

func UpMigrateWithFile(version, dirname, fname string, conn gorec.Querier) error {
	fh, err := os.Open(dirname + "/" + fname)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(fh)
	if err != nil {
		return err
	}

	_, err = conn.Exec(string(data))
	if err != nil {
		return err
	}

	conn.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version)

	return nil
}

func DownMigrateWithFile() {

}

func PerformUpMigration(db *sql.DB, identifier string, migration string) {

}

func PerformDownMigration(db *sql.DB, identifier string, migration string) {

}
