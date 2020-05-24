package main

import (
	"os"
)

type GorecConfig struct {
	Model string
	TableName string
	Package string
	PrimaryKey string
	RawPrefix string
	InternalPrefix string
}

func NewGorecConfig() GorecConfig {
	// var curFile := os.Getenv("GOFILE")
	// var curLine := os.Getenv("GOLINE")
	// Deduce the model from the file

	return GorecConfig{
		PrimaryKey: "id",
		Package: os.Getenv("GOPACKAGE"),
		RawPrefix: "Raw",
		InternalPrefix: "Internal",
	}
}
