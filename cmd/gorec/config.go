package main

type GorecConfig struct {
	Model string
	TableName string
	PrimaryKey string
}

func NewGorecConfig() {
	return GorecConfig{
		PrimaryKey: "id",
	}
}
