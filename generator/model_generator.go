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


	// Write file header
	fmt.Fprintf(fh, `package `+cfg.Package+`

// Code generated by GoRecord.  DO NOT EDIT.

import (
	"database/sql"
	"errors"
	"github.com/johnnyb/gorecord/gorec"
)
`)


	// Load Column Data and Prepare
	if cfg.TableName == "" {
		cfg.TableName = inflect.Pluralize(inflect.Underscore(cfg.Model))
	}
	columnInfo := LoadColumnData(db, cfg, cfg.TableName)

	packages := map[string]bool{}
	allDbNames := []string{}
	allDbNamesNoPk := []string{}
	allStructPointers := []string{}
	allStructValues := []string{}
	allStructValuesNoPk := []string{}
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
		} else {
			allDbNamesNoPk = append(allDbNamesNoPk, ctype.DbName)
			allStructValuesNoPk = append(allStructValuesNoPk, "rec." + ctype.StructName)
		}
	}

	allDbNamesNoPkStr := strings.Join(allDbNamesNoPk, ", ")
	allStructValuesNoPkStr := strings.Join(allStructValuesNoPk, ", ")
	allDbNamesStr := strings.Join(allDbNames, ", ")
	allStructValuesStr := strings.Join(allStructValues, ", ")
	allStructPointersStr := strings.Join(allStructPointers, ", ")
	setDbValuesStr := strings.Join(setDbValues, ", ")


	// Write struct
	fmt.Fprintf(fh, `
type `+cfg.Model+`Record struct {
	`+cfg.InternalPrefix+`IsSaved bool
`)

	for _, cdata := range columnInfo {
		fmt.Fprintf(fh, "\t%s %s\n", cdata.StructName, cdata.ColumnType)
	}
	fmt.Fprintf(fh, "}\n\n")

	// Write Getters/Setters
	for _, cdata := range columnInfo {
		cfg.WriteMethod(fh, cdata.FuncName, "() "+cdata.ColumnType, "\treturn rec."+cdata.StructName+"\n")
		cfg.WriteMethod(fh, "Set"+cdata.FuncName, "(val "+cdata.ColumnType+")", "\trec."+cdata.StructName+" = val\n")
	}

	// Write standard functions
	ctxFunc := fmt.Sprintf("%sGlobalTransactionContext", cfg.Model)
	cfg.WriteFunc(fh, "GlobalConnection", "() *sql.DB", "\treturn gorec.GlobalConnection\n")
	cfg.WriteFunc(fh, "GlobalTransactionContext", "() gorec.Querier", "\treturn gorec.GlobalTransactionContext\n")
	cfg.WriteFunc(fh, "New", fmt.Sprintf("() *%s", cfg.Model), fmt.Sprintf("\trec := %s{}\n\trec.InitializeNew()\n\trec.AutoGeneratePrimaryKey()\n\treturn &rec\n", cfg.Model))
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
	} else if (keyColumn.ColumnType == "int32" || keyColumn.ColumnType == "int64") {
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
	valuePlaceholdersNoPkStr := strings.Join(valuePlaceholders[:(len(valuePlaceholders) - 1)], ", ")
	setPrimaryKeyPlaceholder := "$" + fmt.Sprintf("%d", (len(setDbValues) + 1))
	cfg.WriteMethod(fh, "Save", "() error", `
	var err error = nil
	conn := `+ctxFunc+`()
	if rec.IsSaved() {
		_, err = conn.Exec("UPDATE `+cfg.TableName+` SET `+setDbValuesStr+` WHERE `+cfg.PrimaryKey+` = `+setPrimaryKeyPlaceholder+`", `+allStructValuesStr+`, rec.`+keyColumn.StructName+`)
		return err
	} else {
		if !rec.`+cfg.InternalPrefix+`PrimaryKeyIsPresent() {
			// Let the DB generate the PK
			results, err := conn.Query("INSERT INTO `+cfg.TableName+` (`+allDbNamesNoPkStr+`) VALUES (`+valuePlaceholdersNoPkStr+`) RETURNING id", `+allStructValuesNoPkStr+`)
			if err != nil {
				return err
			}
			if !results.Next() {
				return errors.New("Insert did not return ID")
			}
			results.Scan(&rec.`+keyColumn.StructName+`)
		} else {
			// Insert it ourselves
			_, err := conn.Exec("INSERT INTO `+cfg.TableName+` (`+allDbNamesStr+`) VALUES (`+valuePlaceholdersStr+`)", `+allStructValuesStr+`)
			if err != nil {
				return err
			}
		}
		rec.`+cfg.InternalPrefix+`IsSaved = true
		return nil
	}
`)
}
