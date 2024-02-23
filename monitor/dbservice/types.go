package dbservice

import (
	"database/sql"
)

// TYPES
type MySqlConfig struct {
	User   string
	Passwd string
	Net    string
	Addr   string
	DBName string
}

type DBManager struct {
	connection *sql.DB
}
