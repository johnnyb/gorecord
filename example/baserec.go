package main

type BaseRecordInterface interface {
	PrimaryKeyField() string
	DefaultForeignKeyField() string
}

type BaseRecord struct {
	recordData map[string] interface{} 
	connection *Connection
}

type Connection struct {
}


func GetDefaultConnection() *Connection {
	return nil
}

func (rec *BaseRecord) Connection() {
	if rec.connection == nil {
		rec.connection = GetDefaultConnection()
	}
}

func (rec *BaseRecord) RecordData() map[string] interface{} {
	return rec.recordData;
}

func (rec *BaseRecord) PrimaryKeyField() string {
	return "id"
}

func (rec *BaseRecord) DefaultForeignKeyField() string {
	return "foreign_key_id"
}

func (rec *BaseRecord) HasManyImplementation(dest *BaseRecord) {
	
}
