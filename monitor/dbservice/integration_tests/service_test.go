package integration_tests

import (
	"database/sql"
	"testing"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	db "github.com/hoodnoah/eve_market/monitor/dbservice"
)

func TestNewMySqlDBService(t *testing.T) {
	// configure connection to the test db
	config := mysql.Config{
		User:                 "testuser",
		Passwd:               "password",
		Net:                  "tcp",
		Addr:                 "localhost:3306",
		DBName:               "dbservice_test",
		AllowNativePasswords: true,
	}

	// initialize the service, which will create the tables
	dbservice, err := db.NewMySqlDBService(&config)
	if err != nil {
		t.Fatalf("Failed to create MySqlDBService: %v", err)
	}
	defer dbservice.Close()

	// open a separate connection to test values
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		t.Fatalf("Failed to connect independently to the database: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec("SELECT 1 FROM market_data LIMIT 1"); err != nil {
		t.Errorf("market_data table does not exist or is not accessible: %v", err)
	}

	if _, err := db.Exec("SELECT 1 FROM completed_dates LIMIT 1"); err != nil {
		t.Errorf("completed_dates table does not exist or is not accessible: %v", err)
	}
}
