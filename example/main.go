package main

import (
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/stdlib"
	"database/sql"
	"fmt"
)

func main() {
	db, err := sql.Open("pgx", "user=postgres password=secret host=localhost port=5432 database=gotest sslmode=disable")
	if err != nil {
		panic(err)
	}

	PersonRecordSharedConnection = db
	defer db.Close()

	p := Person{}
	p.SetName("Fred")
	p.SetDescription("Bah")
	err = p.Save()
	if err != nil {
		panic(err)
	}

	p2, err := FindPerson("c20d2036-9cd0-44ed-9b9e-3dfd683d9347")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found person: %s\n", p2.Name());
	p2.SetName(p2.Name() + "Again")
	err = p2.Save()
	if err != nil {
		panic(err)
	}
}
