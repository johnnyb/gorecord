package main

type GorecConfig struct {
	Model string
	TableName string
	PrimaryKey string
}

func NewGorecConfig() GorecConfig {
	return GorecConfig{
		PrimaryKey: "id",
	}
}
