package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lewjian/dbHelper/config"
	"github.com/lewjian/dbHelper/db"
	"time"
)

type MySQL struct {
	User     string
	Host     string
	Password string
	Port     int
	Charset  string
	Db       *sql.DB
	Database string
}

func NewMySQL(user string, host string, password string, port int, charset string, database string) *MySQL {
	return &MySQL{User: user, Host: host, Password: password, Port: port, Charset: charset, Database: database}
}

func (ms *MySQL) Connect() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=%s", ms.User, ms.Password, ms.Host,
		ms.Port, ms.Charset)
	h, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	ms.Db = h
	return nil
}

func (ms *MySQL) Name() string {
	return "mysql"
}

func (ms *MySQL) ListDatabases() ([]string, error) {
	rows, err := ms.Db.Query(`show databases `)
	if err != nil {
		return nil, err
	}
	var mss []string
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			return nil, err
		}
		mss = append(mss, name)
	}
	return mss, nil
}

func (ms *MySQL) ListTables() ([]db.TableInfo, error) {
	rows, err := ms.Db.Query("select TABLE_NAME, ENGINE, TABLE_ROWS, DATA_LENGTH, INDEX_LENGTH, AUTO_INCREMENT, TABLE_COMMENT, CREATE_TIME from information_schema.`TABLES` where TABLE_SCHEMA = ?", ms.Database)
	if err != nil {
		return nil, err
	}
	var tables []db.TableInfo
	for rows.Next() {
		var t db.TableInfo
		var autoIncrement sql.NullInt64
		var createTime sql.NullString
		if err = rows.Scan(&t.Name, &t.Engine, &t.Rows, &t.DataLength, &t.IndexLength, &autoIncrement,
			&t.Comment, &createTime); err != nil {
			return nil, err
		}
		if autoIncrement.Valid {
			t.AutoIncrement = autoIncrement.Int64
		}
		if createTime.Valid {
			t.CreateTime, _ = time.ParseInLocation(config.DefaultDateTimeFormatTpl, createTime.String, time.Local)
		}
		tables = append(tables, t)
	}
	return tables, nil
}

func (ms *MySQL) ListColumns(tableName string) ([]db.ColumnInfo, error) {
	rows, err := ms.Db.Query("select COLUMN_NAME,IS_NULLABLE,DATA_TYPE,`COLUMN_COMMENT`,COLUMN_KEY,COLUMN_TYPE  from information_schema.`COLUMNS` where TABLE_NAME = ? and TABLE_SCHEMA = ?",
		tableName, ms.Database)
	if err != nil {
		return nil, err
	}
	var columns []db.ColumnInfo
	for rows.Next() {
		var t db.ColumnInfo
		var isNullable string
		var comment, key sql.NullString
		if err = rows.Scan(&t.Name, &isNullable, &t.DataType, &comment, &key, &t.ColumnType); err != nil {
			return nil, err
		}
		if isNullable == "YES" {
			t.IsNullable = true
		}
		if comment.Valid {
			t.Comment = comment.String
		}
		if key.Valid && key.String == "PRI" {
			t.IsPrimary = true
		}
		columns = append(columns, t)
	}
	return columns, nil
}
