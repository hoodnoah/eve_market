package dbservice

import (
	"database/sql"
	"sync"
	"time"

	"github.com/hoodnoah/eve_market/monitor/idcache"
	"github.com/hoodnoah/eve_market/monitor/logger"
	"github.com/hoodnoah/eve_market/monitor/parser"
)

// TYPES
type MySqlConfig struct {
	User   string
	Passwd string
	Net    string
	Addr   string
	DBName string
}

type DBManager struct {
	connection   *sql.DB
	logger       logger.ILogger
	input        chan *parser.MarketDay
	output       chan time.Time
	numWorkers   uint
	idCache      idcache.IIDcache
	idCacheMutex sync.Mutex
	mutex        sync.Mutex
}

type InsertID struct {
	ID    int
	Value string
}
