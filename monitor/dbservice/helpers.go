package dbservice

import (
	"database/sql"
	"fmt"
	"strings"

	mds "github.com/hoodnoah/eve_market/monitor/parser"
)

// setup the tables
func bootstrapTables(connection *sql.DB) error {
	tx, err := connection.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec(completedDatesTableTemplate); err != nil {
		return err
	}
	if _, err = tx.Exec(marketDataTableTemplate); err != nil {
		return err
	}

	return tx.Commit()
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
