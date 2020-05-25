package main

import (
	"flag"
	"github.com/johnnyb/gorecord/gorec"
	"github.com/johnnyb/gorecord/generator"

	// NOTE - all supported databases must be listed here
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Config struct {
	Action string
}

func main() {
	cfg := generator.NewConfig()
	cmdCfg := Config{}
	cmdCfg.Action = "model"
	parseFlags(&cmdCfg, &cfg)
	db, err := gorec.AutoConnect()
	if err != nil {
		panic(err)
	}

	// Probably should separate out command config from generator config

	// Open cfg.OriginalFile
	// Look for functions
	// Add them to SkipFunctions

	generator.GenerateModelFile(db, cfg)
}

func parseFlags(cmdCfg *Config, cfg *generator.Config) {
	flag.StringVar(&cmdCfg.Action, "action", cmdCfg.Action, "What action to perform - model (default), has_many, has_one, or belongs_to")
	flag.StringVar(&cfg.Model, "model", cfg.Model, "The name of the model to generate")
	flag.StringVar(&cfg.TableName, "table", cfg.TableName, "The name of the table for the model")
	flag.StringVar(&cfg.Package, "pkg", cfg.Package, "The name of the package to use")
	flag.Var((*AppendSliceValue)(&cfg.SkipFunctions), "skipfunction", "One per argument of functions to skip generating")
	flag.Parse()
}
