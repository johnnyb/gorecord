package generator

import (
	"database/sql"
	"fmt"
	"github.com/johnnyb/gorecord/inflect"
	"io"
)

func GenerateBelongsToFunc(db *sql.DB, cfg Config) {
	fname := fmt.Sprintf("%s.impl.go", inflect.Underscore(cfg.Model))

	err := WithFileAppend(fname, func(fh io.Writer) {
		WriteBelongsTo(fh, db, cfg)
	})

	if err != nil {
		panic(err)
	}
}

func WriteBelongsTo(fh io.Writer, db *sql.DB, cfg Config) {
	relationship := cfg.Relationship
	targetModel := cfg.TargetModel
	// targetTable := inflect.Pluralize(inflect.Underscore(targetModel))
	if targetModel == "" {
		targetModel = inflect.Singularize(relationship)
	}
	targetColumnName := cfg.ForeignKey
	if targetColumnName == "" {
		targetColumnName = inflect.Underscore(targetModel) + "_id"
	}

	fmt.Fprintf(fh, `func (rec *`+cfg.Model+`) `+targetModel+`() (*`+targetModel+`, error) {
	recs, err := `+targetModel+`QuerySimple("where id = $1", rec.`+ColumnToStructMember(targetColumnName)+`)
	if err != nil {
		return nil, err
	}
	if len(recs) == 0 {
		return nil, nil
	}
	return recs[0], nil
}
`)
}
