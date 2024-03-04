package dbservice

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/hoodnoah/eve_market/monitor/idcache"
	"github.com/hoodnoah/eve_market/monitor/logger"
	"github.com/hoodnoah/eve_market/monitor/parser"
	"github.com/hoodnoah/eve_market/monitor/util"
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

	// setup an idCache
	idCacheClient := idcache.NewRateLimitedClient(logger, 1, 3, 2.5)
	idCache := idcache.NewIDCache(logger, idCacheClient)
	knownRegionIds, err := fetchKnownIDS(conn, idcache.RegionID)
	if err != nil {
		return nil, err
	}
	knownTypeIds, err := fetchKnownIDS(conn, idcache.TypeID)
	if err != nil {
		return nil, err
	}
	idCache.SetKnownIDs(knownRegionIds)
	idCache.SetKnownIDs(knownTypeIds)

	return &DBManager{
		connection:   conn,
		logger:       logger,
		input:        inputChannel,
		output:       outputChannel,
		numWorkers:   numWorkers,
		idCache:      idCache,
		idCacheMutex: sync.Mutex{},
		mutex:        sync.Mutex{},
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

func (dm *DBManager) insertNewRegionAndTypeIds(date *parser.MarketDay) error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	dm.logger.Info(fmt.Sprintf("fetching id data for %s...", date.Date.Format(time.DateOnly)))

	// enumerate all the day's ids
	// using maps as sets to prevent duplication
	regionIdsToLabel := []int{}
	typeIdsToLabel := []int{}

	for _, day := range date.Records {
		regionIdsToLabel = append(regionIdsToLabel, int(day.RegionID))
		typeIdsToLabel = append(typeIdsToLabel, int(day.TypeID))
	}

	// label them
	labeledRegionIds, err := (dm.idCache).LabelMany(regionIdsToLabel)
	if err != nil {
		return err
	}

	labeledTypeIds, err := (dm.idCache).LabelMany(typeIdsToLabel)
	if err != nil {
		return err
	}

	// convert to list of structs for chunking
	regionIDSInsertList := make([]InsertID, 0, len(labeledRegionIds))
	typeIDSInsertList := make([]InsertID, 0, len(labeledTypeIds))
	for k, v := range labeledRegionIds {
		regionIDSInsertList = append(regionIDSInsertList, InsertID{ID: k, Value: v})
	}
	for k, v := range labeledTypeIds {
		typeIDSInsertList = append(typeIDSInsertList, InsertID{ID: k, Value: v})
	}

	// add them to the database
	tx, err := dm.connection.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, chunk := range util.ChunkSlice(regionIDSInsertList, MAXCHUNKSIZE*8/2) {
		res := prepIDInsertQuery(idcache.RegionID, chunk)
		_, err := tx.Exec(res.Query, res.Args...)
		if err != nil {
			return err
		}
	}

	for _, chunk := range util.ChunkSlice(typeIDSInsertList, MAXCHUNKSIZE*8/2) {
		res := prepIDInsertQuery(idcache.TypeID, chunk)
		_, err := tx.Exec(res.Query, res.Args...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
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
	chunks := util.ChunkSlice(day.Records, MAXCHUNKSIZE)
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

				// handle new ids first
				err := dm.insertNewRegionAndTypeIds(marketDay)
				if err != nil {
					errMsg := fmt.Sprintf("Failed to insert the ids for %s: %s", marketDay.Date.Format(time.DateOnly), err)
					dm.logger.Error(errMsg)
					continue
				}

				// then insert the day's data
				res, err := dm.InsertMarketDay(marketDay)
				if err != nil {
					errMsg := fmt.Sprintf("Failed to insert the data for %s: %s", marketDay.Date.Format(time.DateOnly), err)
					dm.logger.Error(errMsg)
				} else {
					dm.output <- res
				}
			}
		}()
	}
}
