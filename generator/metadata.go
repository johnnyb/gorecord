package generator

import (
	"database/sql"
	"fmt"
	"github.com/johnnyb/gorecord/inflect"
)

type ColumnData struct {
	DbName        string
	StructName    string
	FuncName      string
	ColumnType    string
	ColumnPackage string
	Nullable      bool
}

func LoadColumnData(db *sql.DB, cfg Config, tableName string) []ColumnData {
	columnInfo := []ColumnData{}

	qstring := fmt.Sprintf("SELECT * FROM %s LIMIT 1", tableName)
	rows, err := db.Query(qstring)
	panicIfError(err)
	defer rows.Close()

	ctypes, err := rows.ColumnTypes()
	panicIfError(err)

	if len(ctypes) == 0 {
		panic(fmt.Sprintf("No columns found in table '%s'.  Some database drivers currently require a row to be present in the database to introspect them.", tableName))
	}

	for _, ctype := range ctypes {
		name := ctype.Name()
		nullable, _ := ctype.Nullable()
		// precision, scale, ok := ctype.DecimalSize()
		// dbtype := ctype.DatabaseTypeName()
		tp := ctype.ScanType()
		// tpName := tp.Name()
		tpPackage := tp.PkgPath()
		tpPartial := tp.String()

		if nullable {
			tpPartial = "*" + tpPartial
		}

		funcName := inflect.Camelize(name)
		structName := cfg.RawPrefix + funcName

		cdata := ColumnData{
			DbName:        name,
			StructName:    structName,
			FuncName:      funcName,
			ColumnType:    tpPartial,
			ColumnPackage: tpPackage,
			Nullable:      nullable,
		}

		columnInfo = append(columnInfo, cdata)
	}

	return columnInfo
}
