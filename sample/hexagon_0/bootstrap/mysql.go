package bootstrap

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// DDL:
// CREATE TABLE users (
//   id   INT          NOT NULL AUTO_INCREMENT PRIMARY KEY,
//   name VARCHAR(255) NOT NULL
// );

func NewMysql(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
