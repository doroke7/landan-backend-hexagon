package bootstrap

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DDL:
// CREATE TABLE users (
//   id   INT          NOT NULL AUTO_INCREMENT PRIMARY KEY,
//   name VARCHAR(255) NOT NULL
// );

func NewMysql() (*gorm.DB, error) {
	var sHost, sPort string
	if len(CONFIG.DATABASE.WRITE.HOSTS) > 0 {
		sHost = CONFIG.DATABASE.WRITE.HOSTS[0]
	}
	if len(CONFIG.DATABASE.WRITE.PORTS) > 0 {
		sPort = CONFIG.DATABASE.WRITE.PORTS[0]
	}

	sDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		CONFIG.DATABASE.USER,
		CONFIG.DATABASE.PASSWORD,
		sHost,
		sPort,
		CONFIG.DATABASE.NAME,
		CONFIG.DATABASE.CHARSET,
	)

	oDB, err := gorm.Open(mysql.Open(sDSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	fmt.Println("[INFO] MySQL 連線. ", sDSN)
	oSqlDB, err := oDB.DB()
	if err != nil {
		return nil, err
	}
	oSqlDB.SetMaxIdleConns(CONFIG.DATABASE.MAX_IDLE_CONNECTIONS)

	return oDB, nil
}
