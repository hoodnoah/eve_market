package dbservice

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/hoodnoah/eve_market/monitor/logger"
	"github.com/hoodnoah/eve_market/monitor/parser"
)

// constructor for
func NewDBManager(config *mysql.Config, logger logger.ILogger, inputChannel chan *parser.MarketDay, outputChannel chan time.Time, numWorkers uint) (*DBManager, error) {
	if logger == nil {
		panic("Could not initialize a new DBManager; logger must not be nil")
	}

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
		logger:     logger,
		input:      inputChannel,
		output:     outputChannel,
		numWorkers: numWorkers,
		mutex:      sync.Mutex{},
	}, nil
}

func (dm *DBManager) Close() error {
	return dm.connection.Close()
}

// gets a list of all dates successfully inserted
func (dm *DBManager) GetCompletedDates() ([]time.Time, error) {
	dates := make([]time.Time, 0)
	statement, err := dm.connection.Prepare("SELECT DISTINCT date FROM completed_dates ORDER BY date ASC")
	if err != nil {
		return nil, err
	}

	rows, err := statement.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var date time.Time
		rows.Scan(&date)
		dates = append(dates, date)
	}

	return dates, nil
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
	chunks := chunkSlice(day.Records, MAXCHUNKSIZE)
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

func (dm *DBManager) Start() {
	for range dm.numWorkers {
		go func() {
			for marketDay := range dm.input {
				res, err := dm.InsertMarketDay(marketDay)
				if err != nil {
					dm.mutex.Lock()
					errMsg := fmt.Sprintf("Failed to insert the data for %s: %s", marketDay.Date.Format(time.DateOnly), err)
					dm.logger.Error(errMsg)
					dm.mutex.Unlock()
				} else {
					dm.mutex.Lock()
					dm.output <- res
					dm.mutex.Unlock()
				}
			}
		}()
	}
}
