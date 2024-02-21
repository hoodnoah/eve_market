package dbservice

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"

	mds "github.com/hoodnoah/eve_market/monitor/marketdataservice"
)

// CONSTANTS
var marketDataTableTemplate = `
	CREATE TABLE IF NOT EXISTS market_data (
		date DATE NOT NULL,
		region_id INTEGER UNSIGNED NOT NULL,
		type_id INTEGER UNSIGNED NOT NULL,
		average DECIMAL(20, 2) NOT NULL,
		highest DECIMAL(20, 2) NOT NULL,
		lowest DECIMAL(20, 2) NOT NULL,
		volume INTEGER UNSIGNED NOT NULL,
		order_count INTEGER UNSIGNED NOT NULL,

		PRIMARY KEY(date, region_id, type_id)
	);
`

var completedDatesTableTemplate = `
	CREATE TABLE IF NOT EXISTS completed_dates (
		date DATE PRIMARY KEY
	);
`

var insertRecordTemplate = `
	INSERT INTO market_data (date, region_id, type_id, average, highest, lowest, volume, order_count)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE date=date;
`

var insertManyTemplate = `
	INSERT INTO market_data
		(date, region_id, type_id, average, highest, lowest, volume, order_count)
	VALUES
		%s
	ON DUPLICATE KEY UPDATE date=date;
`

const MAXCHUNKSIZE = 2000

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
	statement, err := service.connection.Prepare(insertRecordTemplate)
	if err != nil {
		return err
	}
	_, err = statement.Exec(
		record.Date,
		record.RegionID,
		record.TypeID,
		record.Average,
		record.Highest,
		record.Lowest,
		record.Volume,
		record.OrderCount,
	)
	return err
}

// chunk a slice into slices of at most size n, where
// n is the provided chunkSize
func chunkSlice[T any](slice []T, chunkSize int) [][]T {
	returnValue := make([][]T, 0, 0)
	if len(slice) <= chunkSize {
		return append(returnValue, slice)
	}

	idx := 0
	for idx < len(slice) {
		high := idx + min(chunkSize, len(slice)-idx)
		returnValue = append(returnValue, slice[idx:high])
		idx = high
	}

	return returnValue
}

func prepQueryForChunk(chunk []mds.MarketHistoryCSVRecord) struct {
	Query string
	Args  []interface{}
} {
	numRecords := len(chunk)
	placeholders := make([]string, 0, numRecords)
	args := make([]interface{}, 0, numRecords*8)
	for _, record := range chunk {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?, ?)")
		args = append(args, record.Date, record.RegionID, record.TypeID, record.Average, record.Highest, record.Lowest, record.Volume, record.OrderCount)
	}

	return struct {
		Query string
		Args  []interface{}
	}{
		Query: fmt.Sprintf(insertManyTemplate, strings.Join(placeholders, ", ")),
		Args:  args,
	}
}

// Insert several records at once
func (service *MySqlDBService) InsertMany(records []mds.MarketHistoryCSVRecord) error {
	chunkedRecords := chunkSlice[mds.MarketHistoryCSVRecord](records, MAXCHUNKSIZE)

	tx, err := service.connection.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, chunk := range chunkedRecords {
		res := prepQueryForChunk(chunk)
		_, err := tx.Exec(res.Query, res.Args...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
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
