package main

import (
	"fmt"
	"bufio"
	"io/ioutil"
	"os"
	"database/sql"
)

func generateModel(db *sql.DB, cfg GorecConfig) {
	fname := fmt.Sprintf("%s.impl.go", cfg.TableName)
	f, err := os.Create(fname)
	defer f.Close()

	fh := bufio.NewWriter(f)
	defer fh.Flush()

	qstring := fmt.Sprintf("SELECT * FROM %s LIMIT 1", cfg.TableName)
	rows, err := db.Query(qstring)
	ctypes := rows.ColumnTypes()

	columnNames := []string
	columnTypes := []string
	packages := map[string]bool{}

	for ctype := range ctypes {
		name := ctype.Name()
		nullable, _ := ctype.Nullable()
		precision, scale, ok = ctype.DecimalSize()
		dbtype := ctype.DatabaseTypeName()
		tp := ctype.ScanType()
		tpName := tp.Name()
		tpPackage := tp.PkgPath()
		tpPartial := tp.String()

		if tpPackage != "" {
			packages[tpPackage] := true
		}
		columnNames := columnNames.append(name)
		if nullable {
			columnTypes := columnTypes.append(tpPartial)
		} else {
			columnTypes := columnTypes.append("*" + tpPartial)
		}
	}

	fmt.Fprintf(fh, "package %s\n\nimport (\n\t"database/sql"\n)\n\ntype PersonRecord struct {\n")
	for i := 0; i < len(columnNames); i++ {
		fmt.Fprintf("\t%s %s%s")
	}
	fmt.Printf(fh, "}\n")
}
