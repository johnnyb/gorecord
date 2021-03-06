package main

import (
	"flag"
	"github.com/johnnyb/gorecord/generator"
	"github.com/johnnyb/gorecord/gorec"

	// NOTE - all supported databases must be listed here
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Config struct {
	Action    string
	Directory string
	Name      string
}

func main() {
	cfg := generator.NewConfig()
	cmdCfg := Config{
		Directory: "migrations",
		Name:      "migration",
	}
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

	switch cmdCfg.Action {
	case "HasMany":
		generator.GenerateHasManyFunc(db, cfg)

	case "BelongsTo":
		generator.GenerateBelongsToFunc(db, cfg)

	case "model":
		generator.GenerateModelFile(db, cfg)

	case "migration":
		generator.GenerateMigrationFile(cmdCfg.Directory, cmdCfg.Name)
	}
}

func parseFlags(cmdCfg *Config, cfg *generator.Config) {
	flag.StringVar(&cmdCfg.Action, "action", cmdCfg.Action, "What action to perform - model (default), has_many, has_one, or belongs_to")
	flag.StringVar(&cfg.Model, "model", cfg.Model, "The name of the model to generate")
	flag.StringVar(&cfg.TableName, "table", cfg.TableName, "The name of the table for the model")
	flag.StringVar(&cfg.Package, "pkg", cfg.Package, "The name of the package to use")
	flag.StringVar(&cfg.Relationship, "relationship", cfg.Relationship, "The relationship to generate")
	flag.StringVar(&cfg.TargetModel, "targetmodel", cfg.TargetModel, "The target model to use")
	flag.StringVar(&cfg.ForeignKey, "foreignkey", cfg.ForeignKey, "The foreign key to use")
	flag.StringVar(&cmdCfg.Directory, "directory", cmdCfg.Directory, "The directory for migrations (default is `migrations`)")
	flag.StringVar(&cmdCfg.Name, "named", cmdCfg.Name, "The name of a migration")
	flag.Var((*AppendSliceValue)(&cfg.SkipFunctions), "skipfunction", "One per argument of functions to skip generating")
	flag.Parse()
}
