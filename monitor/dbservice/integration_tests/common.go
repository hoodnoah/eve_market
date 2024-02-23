package integration_tests

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	dbs "github.com/hoodnoah/eve_market/monitor/dbservice"
	mds "github.com/hoodnoah/eve_market/monitor/marketdataservice"
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

func clearTables(conn *sql.DB) error {
	tx, err := conn.Begin()
	if err != nil {
		return fmt.Errorf("Transaction failed: %v", err)
	}
	defer tx.Rollback()

	tables := []string{
		"completed_dates",
		"market_data",
	}
	for _, table := range tables {
		_, err := tx.Exec("DELETE FROM " + table + ";")
		if err != nil {
			return fmt.Errorf("Failed to create truncate query for table %s: %v", table, err)
		}
		_, err = tx.Exec("ALTER TABLE " + table + " AUTO_INCREMENT=1;")
		if err != nil {
			return fmt.Errorf("Failed to reset auto increment: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("Failed to commit transaction: %v", err)
	}

	return nil
}

type TestSetup struct {
	Connection *sql.DB
	DBManager  *dbs.DBManager
	TearDown   func()
}

func Setup() (*TestSetup, error) {
	connection, err := sql.Open("mysql", MySqlConfig.FormatDSN())
	if err != nil {
		return nil, err
	}

	dbManager, err := dbs.NewDBManager(&MySqlConfig)
	if err != nil {
		return nil, err
	}

	tearDown := func() {
		if err := clearTables(connection); err != nil {
			log.Fatalf("Failed to clear tables: %v", err)
		}
		connection.Close()
		dbManager.Close()
	}

	return &TestSetup{
		Connection: connection,
		DBManager:  dbManager,
		TearDown:   tearDown,
	}, nil
}

type Scannable interface {
	Scan(dest ...any) error
}

// scans a row into a record
func ScanRowToRecord[T Scannable](row T) (*mds.MarketHistoryCSVRecord, error) {
	resultRecord := mds.MarketHistoryCSVRecord{}
	if err := row.Scan(
		&resultRecord.Date,
		&resultRecord.RegionID,
		&resultRecord.TypeID,
		&resultRecord.Average,
		&resultRecord.Highest,
		&resultRecord.Lowest,
		&resultRecord.Volume,
		&resultRecord.OrderCount,
	); err != nil {
		return nil, err
	}
	return &resultRecord, nil
}
