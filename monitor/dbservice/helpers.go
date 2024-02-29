package dbservice

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/hoodnoah/eve_market/monitor/idcache"
	mds "github.com/hoodnoah/eve_market/monitor/parser"
)

// setup the tables
func bootstrapTables(connection *sql.DB) error {
	tx, err := connection.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec(regionIDsTableTemplate); err != nil {
		return err
	}
	if _, err = tx.Exec(typeIDsTableTemplate); err != nil {
		return err
	}
	if _, err = tx.Exec(completedDatesTableTemplate); err != nil {
		return err
	}
	if _, err = tx.Exec(marketDataTableTemplate); err != nil {
		return err
	}

	return tx.Commit()
}

// fetch the existing region ids from the database
func fetchKnownIDS(connection *sql.DB, idType idcache.IDType) (*idcache.KnownIDs, error) {
	ids := make(map[int]string, 0)
	var tableName string
	switch idType {
	case idcache.RegionID:
		tableName = "region_id"
	default:
		tableName = "type_id"
	}

	query := fmt.Sprintf("SELECT DISTINCT id, value FROM %s;", tableName)
	statement, err := connection.Prepare(query)

	if err != nil {
		return nil, err
	}
	rows, err := statement.Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id int
		var value string
		err = rows.Scan(&id, &value)
		if err != nil {
			return nil, err
		}
	}

	err = statement.Close()
	if err != nil {
		return nil, err
	}

	return &idcache.KnownIDs{
		Type: idcache.RegionID,
		IDS:  ids,
	}, nil
}

func prepQueryForChunk(dateID uint, chunk []mds.MarketHistoryCSVRecord) struct {
	Query string
	Args  []interface{}
} {
	numRecords := len(chunk)
	placeholders := make([]string, 0, numRecords)
	args := make([]interface{}, 0, numRecords*8)
	for _, record := range chunk {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?, ?)")
		args = append(args, dateID, record.RegionID, record.TypeID, record.Average, record.Highest, record.Lowest, record.Volume, record.OrderCount)
	}

	return struct {
		Query string
		Args  []interface{}
	}{
		Query: fmt.Sprintf(insertManyTemplate, strings.Join(placeholders, ", ")),
		Args:  args,
	}
}

func prepIDInsertQuery(idType idcache.IDType, idValues []InsertID) struct {
	Query string
	Args  []interface{}
} {
	var queryTemplate string
	switch idType {
	case idcache.RegionID:
		queryTemplate = insertRegionIDsTemplate
	default:
		queryTemplate = insertTypeIDsTemplate
	}

	numIds := len(idValues)
	placeholders := make([]string, 0, numIds)
	args := make([]interface{}, 0, numIds*2)

	for _, insertID := range idValues {
		id := insertID.ID
		value := insertID.Value
		placeholders = append(placeholders, "(?, ?)")
		args = append(args, id, value)
	}

	return struct {
		Query string
		Args  []interface{}
	}{
		Query: fmt.Sprintf(queryTemplate, strings.Join(placeholders, ", ")),
		Args:  args,
	}
}
