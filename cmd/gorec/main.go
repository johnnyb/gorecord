package main

import (
	"flag"
	"github.com/johnnyb/gorecord/gorec"

	// NOTE - all supported databases must be listed here
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	cfg := NewGorecConfig()
	parseFlags(&cfg)

	db, err := gorec.AutoConnect()
	panicIfError(err)

	generateModel(db, cfg)
}

func parseFlags(cfg *GorecConfig) {
	flag.StringVar(&cfg.Model, "model", cfg.Model, "The name of the model to generate")
	flag.StringVar(&cfg.TableName, "table", cfg.TableName, "The name of the table for the model")
	flag.StringVar(&cfg.Package, "pkg", cfg.Package, "The name of the package to use")
	flag.Parse()
}
