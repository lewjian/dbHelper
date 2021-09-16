package db

import (
	"time"
)

type DBer interface {
	Name() string
	Connect() error
	ListDatabases() ([]string, error)
	ListTables() ([]TableInfo, error)
	ListColumns(tableName string) ([]ColumnInfo, error)
}

type TableInfo struct {
	Name          string
	Engine        string
	Rows          int64
	DataLength    int64
	IndexLength   int64
	AutoIncrement int64
	Comment       string
	CreateTime    time.Time
	Columns       []ColumnInfo
}

type ColumnInfo struct {
	Name          string
	IsPrimary     bool
	IsNullable    bool
	DataType      string
	IsDefaultNull bool
	DefaultValue  string
	Comment       string
	ColumnType    string
}
