package dbservice

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

func NewDBService(config *ConfigVars) (*DBService, error) {
	cfg := mysql.Config{
		User:                 config.User,
		Passwd:               config.Passwd,
		Net:                  config.Net,
		Addr:                 config.Addr,
		DBName:               config.DBName,
		AllowNativePasswords: true,
	}

	// get db handle
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	return &DBService{
		connection: db,
	}, nil
}
