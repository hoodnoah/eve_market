package dbservice

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	mds "github.com/hoodnoah/eve_market/monitor/marketdataservice"
)

const CREATE_MARKET_DATA_TABLE_QUERY = `
	CREATE TABLE
	IF NOT EXISTS
	market_data (
		date        DATE NOT NULL,
		region_id   INTEGER NOT NULL,
		type_id     INTEGER NOT NULL,
		average     FLOAT NOT NULL,
		highest     FLOAT NOT NULL,
		lowest      FLOAT NOT NULL,
		volume      INTEGER NOT NULL,
		order_count INTEGER NOT NULL,

		PRIMARY KEY (date, region_id, type_id)
	);
`

const INSERT_QUERY_TEMPLATE = `
	INSERT INTO
	market_data
	(date, region_id, type_id, average, highest, lowest, volume, order_count)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE date=date; // intentional no-op; id gets set to its current value
`

// bootstraps the market_data table
func ensureMarketDataTableExists(db *sql.DB) error {
	statement, err := db.Prepare(CREATE_MARKET_DATA_TABLE_QUERY)

	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec()

	if err != nil {
		return err
	}

	return nil
}

// given a record and a db, inserts the record
func insertMarketDataRecord(record *mds.MarketHistoryCSVRecord, db *sql.DB) error {
	statement, err := db.Prepare(INSERT_QUERY_TEMPLATE)

	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(
		record.Date,
		record.RegionID,
		record.TypeID,
		record.Average,
		record.Highest,
		record.Lowest,
		record.Volume,
		record.OrderCount)

	if err != nil {
		return err
	}

	return nil
}

// insert multiple records at once
func insertMarketDataRecords(records []mds.MarketHistoryCSVRecord, db *sql.DB) error {
	queryStringTemplate := `
		INSERT INTO
		market_data
		(date, region_id, type_id, average, highest, lowest, volume, order_count)
		VALUES %s
		ON DUPLICATE KEY UPDATE date=date;
	`

	valuesPlaceholder := "(?, ?, ?, ?, ?, ?, ?, ?)"

	var valuesPlaceholders []string
	valuesArgs := make([]interface{}, 0, len(records)*8) // interface{} to permit mixed-types in the slice since we have date, int, float, etc.

	// populate a list of placeholders and a list of their arguments
	for i := range len(records) {
		valuesPlaceholders = append(valuesPlaceholders, valuesPlaceholder)
		record := records[i]

		valuesArgs = append(valuesArgs, record.Date, record.RegionID, record.TypeID, record.Average, record.Highest, record.Lowest, record.Volume, record.OrderCount)
	}

	queryStringTemplate = fmt.Sprintf(queryStringTemplate, strings.Join(valuesPlaceholders, ","))

	tpl, err := db.Prepare(queryStringTemplate)
	if err != nil {
		return err
	}

	_, err = tpl.Exec(valuesArgs...)

	return err
}

// constructor for DBService struct
func NewDBService(config *ConfigVars) (*DBService, error) {
	cfg := mysql.Config{
		User:                 config.User,
		Passwd:               config.Passwd,
		Net:                  config.Net,
		Addr:                 config.Addr,
		DBName:               config.DBName,
		AllowNativePasswords: true, // SHOULD BE CHANGED ONCE TESTING IS ESTABLISHED
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

	// handy defaults
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// bootstrap tables, etc.
	err = ensureMarketDataTableExists(db)
	if err != nil {
		return nil, err
	}

	return &DBService{
		connection: db,
	}, nil
}

// resource cleanup for DBService struct
func (d *DBService) Close() {
	d.connection.Close()
}

// adds a record to the database
func (d *DBService) AddRecord(record *mds.MarketHistoryCSVRecord) error {
	return insertMarketDataRecord(record, d.connection)
}

// adds multiple records at once to the database
func (d *DBService) BulkAddRecord(records []mds.MarketHistoryCSVRecord) error {
	return insertMarketDataRecords(records, d.connection)
}
