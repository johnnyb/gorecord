package generator

import (
	"database/sql"
	"fmt"
	"os"
	"io"
	"bufio"
	"github.com/johnnyb/gorecord/inflect"
)

func GenerateHasManyFunc(db *sql.DB, cfg Config) {
	fname := fmt.Sprintf("%s.impl.go", inflect.Underscore(cfg.Model))
	f, err := os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}
	if f == nil {
		panic("Unable to open file!")
	}
	defer f.Close()

	fh := bufio.NewWriter(f)
	defer fh.Flush()

	WriteHasMany(fh, db, cfg)
}

func WriteHasMany(fh io.Writer, db *sql.DB, cfg Config) {
	relationship := cfg.Relationship
	targetModel := cfg.TargetModel
	// targetTable := inflect.Pluralize(inflect.Underscore(targetModel))
	if targetModel == "" {
		targetModel = inflect.Singularize(relationship)
	}
	targetColumnName := cfg.ForeignKey
	if targetColumnName == "" {
		targetColumnName = inflect.Underscore(cfg.Model) + "_id"
	}

	fmt.Fprintf(fh, "func (rec *%s) %s() ([]*%s, error) {\n\treturn %sQuerySimple(\"where %s = %s\", rec.PrimaryKey())\n}\n\n", cfg.Model, cfg.Relationship, targetModel, targetModel, targetColumnName, "$1")
}
