package generator

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/johnnyb/gorecord/inflect"
	"reflect"
	"time"
)

type ColumnData struct {
	DbName        string
	StructName    string
	FuncName      string
	ColumnType    string
	ColumnPackage string
	SqlNullable   bool
	Nullable      bool
}

var boolType = reflect.TypeOf(true)
var stringType = reflect.TypeOf("")
var int16Type = reflect.TypeOf(int16(0))
var int32Type = reflect.TypeOf(int32(0))
var int64Type = reflect.TypeOf(int64(0))
var timeType = reflect.TypeOf(time.Time{})
var byteType = reflect.TypeOf(byte(0))
var float64Type = reflect.TypeOf(float64(0))
var uuidType = reflect.TypeOf(uuid.UUID{})

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
		funcName := inflect.Camelize(name)
		structName := cfg.RawPrefix + funcName
		// precision, scale, ok := ctype.DecimalSize()
		dbtype := ctype.DatabaseTypeName()

		cdata := ColumnData{
			DbName:        name,
			StructName:    structName,
			FuncName:      funcName,
			SqlNullable:   true,
			Nullable:      false,
			ColumnPackage: "database/sql",
		}
		tp := ctype.ScanType()

		if dbtype == "UUID" {
			tp = uuidType
		}

		switch tp {
		case uuidType:
			cdata.ColumnType = "uuid.NullUUID"
			cdata.ColumnPackage = "github.com/google/uuid"
		case stringType:
			cdata.ColumnType = "sql.NullString"
		case int16Type:
			cdata.ColumnType = "sql.NullInt16"
		case int32Type:
			cdata.ColumnType = "sql.NullInt32"
		case int64Type:
			cdata.ColumnType = "sql.NullInt64"
		case timeType:
			cdata.ColumnType = "sql.NullTime"
		case boolType:
			cdata.ColumnType = "sql.NullBool"
		case byteType:
			cdata.ColumnType = "sql.NullByte"
		case float64Type:
			cdata.ColumnType = "sql.NullFloat64"
		default:
			// tpName := tp.Name()
			tpPackage := tp.PkgPath()
			tpPartial := tp.String()

			nullable, ok := ctype.Nullable()
			if !ok {
				nullable = true
			}

			if nullable {
				tpPartial = "*" + tpPartial
			}
			cdata.ColumnType = tpPartial
			cdata.ColumnPackage = tpPackage
			cdata.SqlNullable = false
		}
		columnInfo = append(columnInfo, cdata)

	}

	return columnInfo
}
