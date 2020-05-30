package generator

import (
	"fmt"
	"time"
)

func GenerateMigrationFile(dir, name string) {
	t := time.Now()
	version := fmt.Sprintf("%04d%02d%02d%02d%02d%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	fname := fmt.Sprintf("%s/%s_%s.go", dir, version, name)

	WriteFile(fname, `package `+dir+`

import (
	"github.com/johnnyb/gorecord/gorec"
	"github.com/johnnyb/gorecord/migrator"
)

func init() {
	migrator.RegisterMigration(
		"`+version+`",
		func(conn gorec.Querier) error {
			panic("No up migration implemented")
			return nil
		}, 
		func(conn gorec.Querier) error {
			panic("No down migration implemented")
			return nil
		}, 
	)
}
`)
}
