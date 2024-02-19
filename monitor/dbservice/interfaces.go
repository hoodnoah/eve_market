package dbservice

import (
	_ "github.com/go-sql-driver/mysql"
	mds "github.com/hoodnoah/eve_market/monitor/marketdataservice"
)

type IDataType interface {
	Type() string
	ToMySqlString() string
}

// integer type
type Integer struct{}

func (i Integer) Type() string {
	return "Integer"
}

// fixed decimal type
type FixedDecimal struct {
	Scale     int
	Precision int
}

func (d FixedDecimal) Type() string {
	return "FixedDecimal"
}

type IDBService interface {
	InsertOne(*mds.MarketHistoryCSVRecord) error
	InsertMany([]mds.MarketHistoryCSVRecord) error
	QueryOne(string) (*mds.MarketHistoryCSVRecord, error)
	QueryMany(string) ([]mds.MarketHistoryCSVRecord, error)
	Close() error
}
