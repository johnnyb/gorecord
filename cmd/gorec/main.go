package main

import (
	"flag"
	"fmt"
	"github.com/johnnyb/gorecord/generator"
	"github.com/johnnyb/gorecord/gorec"
	"os"

	// NOTE - all supported databases must be listed here
	_ "github.com/jackc/pgx/v4/stdlib"
)

type Config struct {
	Action           string
	Directory        string
	Name             string
	ConnectionString string
}

func main() {
	cfg := generator.NewConfig()
	cmdCfg := Config{
		Directory: "migrations",
		Name:      "migration",
	}
	cmdCfg.Action = "model"
	parseFlags(&cmdCfg, &cfg)
	if cmdCfg.ConnectionString != "" {
		os.Setenv("DB_CONNECTION_STRING", cmdCfg.ConnectionString)
	}

	db, err := gorec.AutoConnect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to process database connection: %s\n", err.Error())
		os.Exit(1)
	}
	err = db.Ping()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %s\n", err.Error())
		os.Exit(1)
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
		if cfg.Model == "" {
			fmt.Fprintf(os.Stderr, "No model specified\n")
			os.Exit(1)
		}
		generator.GenerateModelFile(db, cfg)

	case "migration":
		generator.GenerateMigrationFile(cmdCfg.Directory, cmdCfg.Name)
	}
}

func parseFlags(cmdCfg *Config, cfg *generator.Config) {
	flag.StringVar(&cmdCfg.ConnectionString, "connection-string", cmdCfg.ConnectionString, "Database connection string (key=val key=val)")
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
