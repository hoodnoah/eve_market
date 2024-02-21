package integration_tests

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	dbs "github.com/hoodnoah/eve_market/monitor/dbservice"
)

var MySqlConfig mysql.Config = mysql.Config{
	User:                 "testuser",
	Passwd:               "password",
	Net:                  "tcp",
	Addr:                 "localhost:3306",
	DBName:               "dbservice_test",
	AllowNativePasswords: true,
	ParseTime:            true,
}

const ClearTablesQuery = `
	TRUNCATE TABLE market_data;
	TRUNCATE TABLE completed_dates;
`

func clearTables(context context.Context, conn *sql.DB) error {
	tx, err := conn.BeginTx(context, nil)
	if err != nil {
		return fmt.Errorf("Transaction failed: %v", err)
	}
	defer tx.Rollback()

	tables := []string{
		"market_data",
		"completed_dates",
	}
	for _, table := range tables {
		_, err := tx.ExecContext(context, "TRUNCATE TABLE "+table+";")
		if err != nil {
			return fmt.Errorf("Failed to create truncate query for table %s: %v", table, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("Failed to commit transaction: %v", err)
	}

	return nil
}

type TestSetup struct {
	Connection *sql.DB
	DBService  *dbs.IDBService
	TearDown   func()
}

func Setup() (TestSetup, error) {
	connection, err := sql.Open("mysql", MySqlConfig.FormatDSN())
	if err != nil {
		return TestSetup{}, err
	}

	dbService, err := dbs.NewMySqlDBService(&MySqlConfig)
	if err != nil {
		return TestSetup{}, err
	}

	tearDown := func() {
		context := context.Background()
		if err := clearTables(context, connection); err != nil {
			log.Fatalf("Failed to clear tables: %v", err)
		}
		connection.Close()
		dbService.Close()
	}

	return TestSetup{
		Connection: connection,
		DBService:  &dbService,
		TearDown:   tearDown,
	}, nil
}
