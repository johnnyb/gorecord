package gorec

import (
	"database/sql"
	"os"
	// "fmt"
)

type Functor func()error
type TransactionFunctor func(* sql.Tx) error

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
	GlobalConnection = conn
	currentTransaction = nil
	GlobalTransactionContext = conn
}

// If you want to automagically set the connection, do this
func AutoConnect() (db *sql.DB, err error) {
	db, err = sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_CONNECTION_STRING"))
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
	if(IsInTransaction()) {
		err := f()
		if err != nil {
			currentTransaction.Rollback()
		}
		return err
	} else {
		WithDBTransaction(GlobalConnection, func(tx *sql.Tx)error {
			GlobalTransactionContext = tx
			currentTransaction = tx
			return f()
		})
	}
	return nil
}

// If you want a one-off transaction, use this
func WithDBTransaction(db *sql.DB, f TransactionFunctor) error {
	tx, err := GlobalConnection.Begin()
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

