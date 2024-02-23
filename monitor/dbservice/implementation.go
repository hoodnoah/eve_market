package dbservice

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"

	mds "github.com/hoodnoah/eve_market/monitor/marketdataservice"
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
func (dm *DBManager) InsertMarketDay(day *mds.MarketDay) (time.Time, error) {
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
	chunks := chunkSlice[mds.MarketHistoryCSVRecord](day.Records, MAXCHUNKSIZE)
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

// // Insert several records at once
// func (dm *DBManager) InsertMany(records []mds.MarketHistoryCSVRecord) error {
// 	chunkedRecords := chunkSlice[mds.MarketHistoryCSVRecord](records, MAXCHUNKSIZE)

// 	tx, err := service.connection.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	defer tx.Rollback()

// 	for _, chunk := range chunkedRecords {
// 		res := prepQueryForChunk(chunk)
// 		_, err := tx.Exec(res.Query, res.Args...)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return tx.Commit()
// }
// func (service *MySqlDBService) QueryOne(string) (*mds.MarketHistoryCSVRecord, error) {
// 	return nil, errors.New("QueryOne not implemented")
// }
// func (service *MySqlDBService) QueryMany(string) ([]mds.MarketHistoryCSVRecord, error) {
// 	return nil, errors.New("QueryMany not implemented")
// }

// func (service MySqlDBService) Close() error {
// 	return service.connection.Close()
// }

// // UTIL
// // creates a given table in the database
// func createTable(connection *sql.DB, query string) error {
// 	_, err := connection.Exec(query)

// 	return err
// }

// // constructor for a DBService
// func NewMySqlDBService(config *mysql.Config) (IDBService, error) {
// 	conn, err := sql.Open(
// 		"mysql",
// 		config.FormatDSN(),
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	// ensure valid connection
// 	err = conn.Ping()
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = createTable(conn, marketDataTableTemplate)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = createTable(conn, completedDatesTableTemplate)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &MySqlDBService{
// 		connection: conn,
// 	}, nil
// }
