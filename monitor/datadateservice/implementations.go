package datadateservice

import (
	"errors"
	"fmt"
	"time"
)

const firstYearAvailable = Year(2003)

var earliestDate = time.Date(int(firstYearAvailable), 10, 1, 0, 0, 0, 0, time.UTC)

const filePrefix = "market-history"
const fileExt = ".csv.bz2"

// lists the days in a given year, exluding anything before
// October 1, 2003, the first date with EveRefs data
func enumerateDatesInYear(year Year) []DataDate {
	if year < 2003 {
		return nil
	}
	startDate := time.Date(int(year), 1, 0, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(int(year)+1, 1, 1, 0, 0, 0, 0, time.UTC)

	var dates []DataDate
	for currentDate := startDate; currentDate.Before(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		if !currentDate.Before(earliestDate) {
			dataDate := DataDate{
				Date: currentDate,
				Url:  createFileURL(currentDate),
			}

			dates = append(dates, dataDate)
		}
	}

	return dates
}

// formats a date into its corresponding CSV file URL
func createFileURL(date time.Time) string {
	fileName := date.Format(time.DateOnly)
	return fmt.Sprintf("%d/%s-%s%s", date.Year(), filePrefix, fileName, fileExt)
}

func createDataYear(year Year, currentDate time.Time) (DataYear, error) {
	if year < 2003 {
		return DataYear{}, fmt.Errorf("Expected a year no earlier than 2003, received :%d", year)
	}

	if year > Year(currentDate.Year()) {
		return DataYear{}, fmt.Errorf("Expected a year no later than the current year %d, received :%d", currentDate.Year(), year)
	}

	allDates := enumerateDatesInYear(year)
	var filteredDates []DataDate

	for _, dataDate := range allDates {
		if !dataDate.Date.After(currentDate) {
			date := dataDate.Date
			dataDate := DataDate{
				Date: date,
				Url:  createFileURL(date),
			}

			filteredDates = append(filteredDates, dataDate)
		}
	}

	if len(filteredDates) < 1 {
		return DataYear{}, errors.New("Could not generate any dates")
	}

	return DataYear{
		Year:  year,
		Dates: filteredDates,
	}, nil
}

// generates a complete list of dates between the provided date
// and the earliest date with data
func (d *DataDateService) EnumerateDataYears(currentDate time.Time) []DataYear {
	numYears := currentDate.Year() - earliestDate.Year() + 1

	var dataYears []DataYear

	for i := range numYears {
		currentYear := Year(earliestDate.Year() + i)
		dates := enumerateDatesInYear(currentYear)
		dataYear := DataYear{
			Year:  currentYear,
			Dates: dates,
		}
		dataYears = append(dataYears, dataYear)
	}

	return dataYears
}
