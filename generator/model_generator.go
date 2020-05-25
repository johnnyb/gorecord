package generator

import (
	"fmt"
	"strings"
	"bufio"
	"io"
	// "io/ioutil"
	"os"
	"database/sql"
	"github.com/johnnyb/gorecord/inflect"
)

type ColumnData struct {
	DbName string
	StructName string
	FuncName string
	ColumnType string
	ColumnPackage string
	Nullable bool
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func GenerateModelFile(db *sql.DB, cfg Config) {
	fname := fmt.Sprintf("%s.impl.go", inflect.Underscore(cfg.Model))
	f, err := os.Create(fname)
	panicIfError(err)
	defer f.Close()

	fh := bufio.NewWriter(f)
	defer fh.Flush()

	WriteModel(fh, db, cfg)
}

func WriteModel(fh io.Writer, db *sql.DB, cfg Config) {
	qstring := fmt.Sprintf("SELECT * FROM %s LIMIT 1", cfg.TableName)
	rows, err := db.Query(qstring)
	panicIfError(err)

	ctypes, err := rows.ColumnTypes()
	panicIfError(err)

	columnInfo := []ColumnData{}
	packages := map[string]bool{}
	allDbNames := []string{}
	allStructPointers := []string{}
	var keyColumn ColumnData

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
		if nullable {
			tpPartial = "*" + tpPartial
		}

		funcName := inflect.Camelize(name)
		structName := cfg.RawPrefix + funcName

		allDbNames = append(allDbNames, name)
		allStructPointers = append(allStructPointers, "&rec." + structName)

		cdata := ColumnData{
			DbName: name,
			StructName: structName,
			FuncName: funcName,
			ColumnType: tpPartial,
			ColumnPackage: tpPackage,
			Nullable: nullable,
		}
		if name == cfg.PrimaryKey {
			keyColumn = cdata
		}

		columnInfo = append(columnInfo, cdata)
	}
	allDbNamesStr := strings.Join(allDbNames, ", ")
	allStructPointersStr := strings.Join(allStructPointers, ", ")

	fmt.Fprintf(fh, "package %s\n\nimport (\n\t\"database/sql\"\n\t\"github.com/johnnyb/gorecord/gorec\"\n)\n\ntype %sRecord struct {\n\t%sIsSaved bool\n", cfg.Package, cfg.Model, cfg.InternalPrefix)
	for _, cdata := range columnInfo {
		fmt.Fprintf(fh, "\t%s %s\n", cdata.StructName, cdata.ColumnType)
	}
	fmt.Fprintf(fh, "}\n\n")
	for _, cdata := range columnInfo {
		fmt.Fprintf(fh, "func (rec *%s) %s() %s {\n\treturn rec.%s\n}\n\n", cfg.Model, cdata.FuncName, cdata.ColumnType, cdata.StructName)
		fmt.Fprintf(fh, "func (rec *%s) Set%s(val %s) {\n\trec.%s = val\n}\n\n", cfg.Model, cdata.FuncName, cdata.ColumnType, cdata.StructName)
	}

	ctxFunc := fmt.Sprintf("%sGlobalTransactionContext", cfg.Model)
	cfg.WriteFunc(fh, "GlobalConnection", "() *sql.DB", "\treturn gorec.GlobalConnection\n")
	cfg.WriteFunc(fh, "GlobalTransactionContext", "() gorec.Querier", "\treturn gorec.GlobalTransactionContext\n")
	cfg.WriteFunc(fh, "New", fmt.Sprintf("() *%s", cfg.Model), fmt.Sprintf("\trec := %s{}\n\trec.InitializeNew()\n\treturn &rec\n", cfg.Model))
	cfg.WriteFunc(fh, "Find", fmt.Sprintf("(key %s) (*%s, error)", keyColumn.ColumnType, cfg.Model), fmt.Sprintf("\trows, err := %s().Query(\"SELECT %s FROM %s WHERE %s = $1\", key)\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\tdefer rows.Close()\n\n\trec := %s{\n\t\t%sRecord{\n\t\t\t%sIsSaved: true,\n\t\t},\n\t}\n\tif rows.Next() {\n\t\trows.Scan(%s)\n\t\trec.InitializeExisting()\n\t\treturn &rec, nil\n\t} else {\n\t\treturn nil, nil\n\t}\n", ctxFunc, allDbNamesStr, cfg.TableName, keyColumn.DbName, cfg.Model, cfg.Model, cfg.InternalPrefix, allStructPointersStr))
	cfg.WriteMethod(fh, "InitializeNew", "()", fmt.Sprintf("\trec.%sIsSaved = false\n", cfg.InternalPrefix))
	cfg.WriteMethod(fh, "InitializeExisting", "()", fmt.Sprintf("\trec.%sIsSaved = true\n", cfg.InternalPrefix))
	cfg.WriteMethod(fh, cfg.InternalPrefix + "ScanAllColumns", "(scanner gorec.RowScanner) error", fmt.Sprintf("\terr := scanner.Scan(%s)\n\tif err != nil {\n\t\treturn err\n\t}\n\trec.InitializeExisting()\n\treturn nil\n", allStructPointersStr))
}
