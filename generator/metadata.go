package generator

import (
	"fmt"
	"database/sql"
	"github.com/johnnyb/gorecord/inflect"
)

func LoadColumnData(db *sql.DB, cfg Config, tableName string) []ColumnData {
	columnInfo := []ColumnData{}

	qstring := fmt.Sprintf("SELECT * FROM %s LIMIT 1", tableName)
	rows, err := db.Query(qstring)
	panicIfError(err)
	defer rows.Close()

	ctypes, err := rows.ColumnTypes()
	panicIfError(err)

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
			DbName: name,
			StructName: structName,
			FuncName: funcName,
			ColumnType: tpPartial,
			ColumnPackage: tpPackage,
			Nullable: nullable,
		}

		columnInfo = append(columnInfo, cdata)
	}

	return columnInfo
}


