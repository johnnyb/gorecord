package main

import (
	"flag"
	"github.com/johnnyb/gorecord/gorec"
	"github.com/johnnyb/gorecord/generator"

	// NOTE - all supported databases must be listed here
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	cfg := generator.NewConfig()
	parseFlags(&cfg)
	db, err := gorec.AutoConnect()
	if err != nil {
		panic(err)
	}

	generator.GenerateModelFile(db, cfg)
}

func parseFlags(cfg *generator.Config) {
	flag.StringVar(&cfg.Model, "model", cfg.Model, "The name of the model to generate")
	flag.StringVar(&cfg.TableName, "table", cfg.TableName, "The name of the table for the model")
	flag.StringVar(&cfg.Package, "pkg", cfg.Package, "The name of the package to use")
	flag.Var((*AppendSliceValue)(&cfg.SkipFunctions), "skipfunction", "One per argument of functions to skip generating")
	flag.Parse()
}
