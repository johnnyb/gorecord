package main

import (
	"fmt"
	"bufio"
	// "io/ioutil"
	"os"
	"database/sql"
)

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func generateModel(db *sql.DB, cfg GorecConfig) {
	fname := fmt.Sprintf("%s.impl.go", cfg.TableName)
	f, err := os.Create(fname)
	defer f.Close()

	fh := bufio.NewWriter(f)
	defer fh.Flush()

	qstring := fmt.Sprintf("SELECT * FROM %s LIMIT 1", cfg.TableName)
	rows, err := db.Query(qstring)
	panicIfError(err)

	ctypes, err := rows.ColumnTypes()
	panicIfError(err)

	columnNames := []string{}
	columnTypes := []string{}
	packages := map[string]bool{}

	for _, ctype := range ctypes {
		name := ctype.Name()
		nullable, _ := ctype.Nullable()
		// precision, scale, ok := ctype.DecimalSize()
		// dbtype := ctype.DatabaseTypeName()
		tp := ctype.ScanType()
		// tpName := tp.Name()
		tpPackage := tp.PkgPath()
		tpPartial := tp.String()

		if tpPackage != "" {
			packages[tpPackage] = true
		}
		columnNames = append(columnNames, name)
		if nullable {
			columnTypes = append(columnTypes, tpPartial)
		} else {
			columnTypes = append(columnTypes, "*" + tpPartial)
		}
	}

	fmt.Fprintf(fh, "package %s\n\nimport (\n\t\"database/sql\"\n)\n\ntype PersonRecord struct {\n")
	for i := 0; i < len(columnNames); i++ {
		cname := columnNames[i]
		ctype := columnTypes[i]
		fmt.Fprintf(fh, "\t%s %s\n", cname, ctype)
	}
	fmt.Fprintf(fh, "}\n")
}
