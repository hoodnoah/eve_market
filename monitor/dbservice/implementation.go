package dbservice

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/hoodnoah/eve_market/monitor/parser"
)

// constructor for
func NewDBManager(config *mysql.Config) (*DBManager, error) {
	// instantiate db conn
	conn, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, err
	}

	// setup all tables
	if err := bootstrapTables(conn); err != nil {
		return nil, err
	}

	return &DBManager{
		connection: conn,
	}, nil
}

func (dm *DBManager) Close() error {
	return dm.connection.Close()
}

// gets a list of all dates successfully inserted
func (dm *DBManager) GetCompletedDates() []time.Time {
	return make([]time.Time, 0)
}

// tries to insert an entire market day's data
// fails unless the entire day can be inserted at once
func (dm *DBManager) InsertMarketDay(day *parser.MarketDay) (time.Time, error) {
	tx, err := dm.connection.Begin()
	if err != nil {
		return time.Time{}, err
	}
	defer tx.Rollback()

	// insert the day's date, keep the id
	result, err := tx.Exec(insertCompletedDateTemplate, day.Date)
	if err != nil {
		return time.Time{}, err
	}

	dateId, err := result.LastInsertId()
	if err != nil {
		return time.Time{}, err
	}

	// chunk the records
	chunks := chunkSlice[parser.MarketHistoryCSVRecord](day.Records, MAXCHUNKSIZE)
	for _, chunk := range chunks {
		queryParts := prepQueryForChunk(uint(dateId), chunk)
		_, err := tx.Exec(queryParts.Query, queryParts.Args...)
		if err != nil {
			return time.Time{}, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return time.Time{}, err
	}
	return day.Date, nil
}
