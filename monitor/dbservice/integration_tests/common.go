package integration_tests

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/hoodnoah/eve_market/monitor/dbservice"
	"github.com/hoodnoah/eve_market/monitor/logger"
	"github.com/hoodnoah/eve_market/monitor/parser"
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
		return fmt.Errorf("transaction failed: %v", err)
	}
	defer tx.Rollback()

	tables := []string{
		"completed_dates",
		"market_data",
	}
	for _, table := range tables {
		_, err := tx.Exec("DELETE FROM " + table + ";")
		if err != nil {
			return fmt.Errorf("failed to create truncate query for table %s: %v", table, err)
		}
		_, err = tx.Exec("ALTER TABLE " + table + " AUTO_INCREMENT=1;")
		if err != nil {
			return fmt.Errorf("failed to reset auto increment: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

type TestSetup struct {
	Connection    *sql.DB
	Logger        logger.ILogger
	DBManager     *dbservice.DBManager
	InputChannel  chan *parser.MarketDay
	OutputChannel chan time.Time
	TearDown      func()
}

type DummyLogger struct{}

var _ logger.ILogger = (*DummyLogger)(nil)

func (d *DummyLogger) Debug(_ string) {}
func (d *DummyLogger) Info(_ string)  {}
func (d *DummyLogger) Warn(_ string)  {}
func (d *DummyLogger) Error(_ string) {}
func (d *DummyLogger) Start()         {}

func Setup(inputChannel chan *parser.MarketDay, outputChannel chan time.Time) (*TestSetup, error) {
	connection, err := sql.Open("mysql", MySqlConfig.FormatDSN())
	if err != nil {
		return nil, err
	}

	dummyLogger := &DummyLogger{}

	dbManager, err := dbservice.NewDBManager(&MySqlConfig, dummyLogger, inputChannel, outputChannel, 1)
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
		Connection:    connection,
		DBManager:     dbManager,
		TearDown:      tearDown,
		Logger:        dummyLogger,
		InputChannel:  inputChannel,
		OutputChannel: outputChannel,
	}, nil
}

type Scannable interface {
	Scan(dest ...any) error
}

// scans a row into a record
func ScanRowToRecord[T Scannable](row T) (*parser.MarketHistoryCSVRecord, error) {
	resultRecord := parser.MarketHistoryCSVRecord{}
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
