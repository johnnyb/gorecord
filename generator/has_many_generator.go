package generator

import (
	"database/sql"
	"fmt"
	"github.com/johnnyb/gorecord/inflect"
	"io"
)

func GenerateHasManyFunc(db *sql.DB, cfg Config) {
	fname := fmt.Sprintf("%s.impl.go", inflect.Underscore(cfg.Model))

	err := WithFileAppend(fname, func(fh io.Writer) {
		WriteHasMany(fh, db, cfg)
	})
	if err != nil {
		panic(err)
	}
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
	targetFieldName := inflect.Camelize(targetColumnName)

	fmt.Fprintf(fh, "func (rec *%s) %s() ([]*%s, error) {\n\treturn %sQuerySimple(\"where %s = %s\", rec.PrimaryKey())\n}\n\n", cfg.Model, cfg.Relationship, targetModel, targetModel, targetColumnName, "$1")
	fmt.Fprintf(fh, "func (rec *%s) %sBuild() (*%s) {\n\tnewrec := %sNew()\n\tnewrec.Set%s(rec.PrimaryKey())\n\treturn newrec\n}\n\n", cfg.Model, cfg.Relationship, targetModel, targetModel, targetFieldName)
}
