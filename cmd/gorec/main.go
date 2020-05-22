package main


import (
	"flag"
	"github.com/johnnyb/gorecord/gorec"
)

func main() {
	cfg := NewGorecConfig()
	parseFlags(&cfg)

	db, err := gorec.GetStandardConnection()
	if err != nil {
		panic(err)
	}

	generateModel(db, cfg)
}

func parseFlags(cfg *GorecConfig) {
	flag.StringVar(&cfg.Model, "model", "", "The name of the model to generate")
	flag.StringVar(&cfg.TableName, "table", "", "The name of the table for the model")
	flag.Parse()
}
