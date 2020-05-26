package generator

import (
	"fmt"
	"strings"
	"bufio"
	"io"
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
	if cfg.TableName == "" {
		cfg.TableName = inflect.Pluralize(inflect.Underscore(cfg.Model))
	}
	columnInfo := LoadColumnData(db, cfg, cfg.TableName)

	packages := map[string]bool{}
	allDbNames := []string{}
	allStructPointers := []string{}
	allStructValues := []string{}
	var keyColumn ColumnData
	setDbValues := []string{}

	for cidx, ctype := range columnInfo {
		if ctype.ColumnPackage != "" {
			packages[ctype.ColumnPackage] = true
		}
		allDbNames = append(allDbNames, ctype.DbName)
		allStructPointers = append(allStructPointers, "&rec." + ctype.StructName)
		allStructValues = append(allStructValues, "rec." + ctype.StructName)
		setDbValues = append(setDbValues, ctype.DbName + " = $" + fmt.Sprintf("%d", cidx))
		if ctype.DbName == cfg.PrimaryKey {
			keyColumn = ctype
		}
	}

	allDbNamesStr := strings.Join(allDbNames, ", ")
	allStructPointersStr := strings.Join(allStructPointers, ", ")
	allStructValuesStr := strings.Join(allStructValues, ", ")
	setDbValuesStr := strings.Join(setDbValues, ", ")

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
	cfg.WriteFunc(fh, cfg.InternalPrefix + "DbAttributeNamesString", "() string", fmt.Sprintf("\treturn \"%s\"\n", allDbNamesStr))
	cfg.WriteFunc(fh, "QuerySimpleWithQuerier", fmt.Sprintf("(querier gorec.Querier, clause string, args ...interface{}) ([]*%s, error)", cfg.Model), fmt.Sprintf("\trows, err := querier.Query(\"SELECT %s FROM %s \" + clause, args...)\n\tresults := []*%s{}\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\tfor rows.Next() {\n\t\tnextres := &%s{}\n\t\terr = nextres.%sScanAllColumns(rows)\n\t\tif err != nil {\n\t\t\treturn results, err\n\t\t}\n\t\tresults = append(results, nextres)\n\t}\n\treturn results, nil\n", allDbNamesStr, cfg.TableName, cfg.Model, cfg.Model, cfg.InternalPrefix))
	cfg.WriteFunc(fh, "QuerySimple", fmt.Sprintf("(clause string, args ...interface{}) ([]*%s, error)", cfg.Model), fmt.Sprintf("\treturn %sQuerySimpleWithQuerier(%s(), clause, args...)\n", cfg.Model, ctxFunc))
	cfg.WriteMethod(fh, "Validate", "() bool", "\treturn true\n")
	cfg.WriteMethod(fh, "PrimaryKey", fmt.Sprintf("() %s", keyColumn.ColumnType), fmt.Sprintf("\treturn rec.%s\n", keyColumn.StructName))
	cfg.WriteMethod(fh, "IsSaved", "() bool", fmt.Sprintf("\treturn rec.%sIsSaved\n", cfg.InternalPrefix))
	if keyColumn.ColumnType == "string" {
		cfg.WriteMethod(fh, fmt.Sprintf("%sPrimaryKeyIsPresent", cfg.InternalPrefix), "() bool", fmt.Sprintf("\treturn rec.%s != \"\"\n", keyColumn.StructName))
	} else if keyColumn.ColumnType == "int32" {
		cfg.WriteMethod(fh, fmt.Sprintf("%sPrimaryKeyIsPresent", cfg.InternalPrefix), "() bool", fmt.Sprintf("\treturn rec.%s != 0\n", keyColumn.StructName))
	} else {
		cfg.WriteMethod(fh, fmt.Sprintf("%sPrimaryKeyIsPresent", cfg.InternalPrefix), "() bool", fmt.Sprintf("\treturn rec.%s != nil\n", keyColumn.StructName))
		
	}
	cfg.WriteMethod(fh, "AutoGeneratePrimaryKey", "()", "\t// Do nothing by default (assume the db will do this)\n")
	valuePlaceholders := []string{}
	for cidx, _ := range columnInfo {
		valuePlaceholders = append(valuePlaceholders, "$" + fmt.Sprintf("%d", cidx + 1))
	}
	valuePlaceholdersStr := strings.Join(valuePlaceholders, ", ")
	setPrimaryKeyPlaceholder := "$" + fmt.Sprintf("%d", (len(setDbValues) + 1))
	cfg.WriteMethod(fh, "Save", "() error", fmt.Sprintf("\tconn := %s()\n\tif rec.IsSaved() {\n\t\t_, err := conn.Exec(\"UPDATE %s SET %s WHERE %s = %s\", %s, rec.%s)\n\t\treturn err\n\t} else {\n\t\tif !rec.%sPrimaryKeyIsPresent() {\n\t\t\trec.AutoGeneratePrimaryKey()\n\t\t}\n\t\t_, err := conn.Exec(\"INSERT INTO %s (%s) VALUES (%s)\", %s)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\trec.%sIsSaved = true\n\t\treturn nil\n\t}\n", ctxFunc, cfg.TableName, setDbValuesStr, cfg.PrimaryKey, setPrimaryKeyPlaceholder, allStructValuesStr, keyColumn.StructName, cfg.InternalPrefix, cfg.TableName, allDbNamesStr, valuePlaceholdersStr, allStructValuesStr, cfg.InternalPrefix))
}