package dbservice

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	mds "github.com/hoodnoah/eve_market/monitor/marketdataservice"
)

type IDBService interface {
	AddRecord(record *mds.MarketHistoryCSVRecord) error
	BulkAddRecord(records []mds.MarketHistoryCSVRecord) error
	ExtantDates() []int
	Close()
}

type ConfigVars struct {
	User   string
	Passwd string
	Net    string
	Addr   string
	DBName string
}

type DBService struct {
	connection      *sql.DB
	bootstrapTables func()
	insert          func(record *mds.MarketHistoryCSVRecord, db sql.DB) error
}
