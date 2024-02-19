package dbservice

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"

	mds "github.com/hoodnoah/eve_market/monitor/marketdataservice"
)

// CONSTANTS
var marketDataTableTemplate = `
	CREATE TABLE IF NOT EXISTS market_data (
		date DATE,
		region_id INTEGER,
		type_id INTEGER,
		average DECIMAL(10, 2),
		highest DECIMAL(10, 2),
		lowest DECIMAL(10, 2),
		volume INTEGER,
		order_count INTEGER,

		PRIMARY KEY(date, region_id, type_id)
	);
`

var completedDatesTableTemplate = `
	CREATE TABLE IF NOT EXISTS completed_dates (
		date DATE PRIMARY KEY
	);
`

// TYPES
type MySqlConfig struct {
	User   string
	Passwd string
	Net    string
	Addr   string
	DBName string
}

type MySqlDBService struct {
	connection *sql.DB
}

// adds a single record to the database
func (service *MySqlDBService) InsertOne(record *mds.MarketHistoryCSVRecord) error {
	return errors.New("InsertOne not implemented")
}
func (service *MySqlDBService) InsertMany(records []mds.MarketHistoryCSVRecord) error {
	return errors.New("InsertMany not implemented")
}
func (service *MySqlDBService) QueryOne(string) (*mds.MarketHistoryCSVRecord, error) {
	return nil, errors.New("QueryOne not implemented")
}
func (service *MySqlDBService) QueryMany(string) ([]mds.MarketHistoryCSVRecord, error) {
	return nil, errors.New("QueryMany not implemented")
}

func (service MySqlDBService) Close() error {
	return service.connection.Close()
}

// UTIL
// creates a given table in the database
func createTable(connection *sql.DB, query string) error {
	_, err := connection.Exec(query)

	return err
}

// constructor for a DBService
func NewMySqlDBService(config *mysql.Config) (IDBService, error) {
	conn, err := sql.Open(
		"mysql",
		config.FormatDSN(),
	)

	if err != nil {
		return nil, err
	}

	// ensure valid connection
	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	err = createTable(conn, marketDataTableTemplate)
	if err != nil {
		return nil, err
	}

	err = createTable(conn, completedDatesTableTemplate)
	if err != nil {
		return nil, err
	}

	return &MySqlDBService{
		connection: conn,
	}, nil
}
