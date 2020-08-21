package gorec

import (
	"database/sql"
	"os"
	// "fmt"
)

type Functor func() error
type TransactionFunctor func(*sql.Tx) error

type Querier interface {
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var currentTransaction *sql.Tx

// This should be the default thing that you
// run queries on.  It might be a transaction
// or it might be the bare *sql.DB.
var GlobalTransactionContext Querier

// If you need the connection itself, here it is.
var GlobalConnection *sql.DB

// If you want to set the connection manually, do this
func SetConnection(conn *sql.DB) {
	currentTransaction = nil

	GlobalConnection = conn
	GlobalTransactionContext = conn
}

func BuildConnectionString(info map[string]string) string {
	cstr := ""
	first := true
	for k, v := range info {
		if !first {
			cstr += " "
		}
		if v != "" {
			cstr += k + "=" + v
			first = false
		}
	}

	return cstr
}

func GetConnectionDriver() string {
	driver := os.Getenv("DB_DRIVER")

	if driver == "" {
		driver = "pgx" // default
	}

	return driver
}

// There are multiple ways to infer a connection string.
func GetConnectionString(testing bool) string {
	var testprefix string
	if testing {
		testprefix = "TEST_"
	}
	cstr := os.Getenv("DB_" + testprefix + "CONNECTION_STRING")
	if cstr == "" {
		dbname := os.Getenv("RDS_" + testprefix + "DB_NAME")
		if dbname != "" {
			host := os.Getenv("RDS_" + testprefix + "HOSTNAME")
			port := os.Getenv("RDS_" + testprefix + "PORT")
			user := os.Getenv("RDS_" + testprefix + "USERNAME")
			pass := os.Getenv("RDS" + testprefix + "PASSWORD")
			cstr := BuildConnectionString(map[string]string{
				"user":     user,
				"dbname":   dbname,
				"host":     host,
				"password": pass,
				"port":     port,
			})
		}
	}

	return cstr
}

// If you want to automagically set the connection, do this
func AutoConnect() (db *sql.DB, err error) {
	db, err = sql.Open(GetConnectionDriver(), GetConnectionString(false))
	if err == nil {
		SetConnection(db)
	}
	return
}

func AutoConnectForTesting() (db *sql.DB, err error) {
	db, err = sql.Open(GetConnectionDriver(), GetConnectionString(true))
	if err == nil {
		SetConnection(db)
	}
	return
}

// Are we currently in a transaction?
func IsInGlobalTransaction() bool {
	return GlobalTransactionContext == GlobalConnection
}

// Execute a function inside a transaction.
// If an error is returned, it rolls back.
// Note that if there is already a transaction
// in progress, this will happen (INCLUDING
// the rollback) in the existing transaction.
func WithGlobalTransaction(f Functor) error {
	if IsInGlobalTransaction() {
		err := f()
		if err != nil {
			currentTransaction.Rollback()
		}
		return err
	} else {
		WithDBTransaction(GlobalConnection, func(tx *sql.Tx) error {
			GlobalTransactionContext = tx
			currentTransaction = tx
			return f()
		})
	}
	return nil
}

// If you want a one-off transaction, use this
func WithDBTransaction(db *sql.DB, f TransactionFunctor) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	err = f(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()

	return err
}
