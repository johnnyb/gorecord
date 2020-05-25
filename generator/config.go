package generator

import (
	"io"
	"os"
)

type Config struct {
	Model string
	TableName string
	Package string
	PrimaryKey string
	RawPrefix string
	InternalPrefix string
	SkipFunctions []string
	OriginalFile string
}

func NewConfig() Config {
	// var curFile := os.Getenv("GOFILE")
	// var curLine := os.Getenv("GOLINE")
	// Deduce the model from the file

	return Config{
		PrimaryKey: "id",
		Package: os.Getenv("GOPACKAGE"),
		RawPrefix: "Raw",
		InternalPrefix: "Internal",
		OriginalFile: os.Getenv("GOFILE"),
	}
}

func (cfg *Config) ShouldSkipFunc(name string) bool {
	for _, f := range cfg.SkipFunctions {
		if f == name {
			return true
		}
	}
	return false
}

func (cfg *Config) WriteMethod(fh io.Writer, name string, signature string, body string) {
	if !cfg.ShouldSkipFunc(name) {
		fh.Write([]byte("func (rec *" + cfg.Model + ") " + name + signature + " {\n" + body + "}\n\n"))
	}
}

func (cfg *Config) WriteFunc(fh io.Writer, name string,signature string, body string) {
	name = cfg.Model + name
	if !cfg.ShouldSkipFunc(name) {
		fh.Write([]byte("func " + name + signature + " {\n" + body + "}\n\n"))
	}
}
