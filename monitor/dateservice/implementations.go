package dateservice

import (
	"time"
)

const firstYearAvailable = Year(2003)

var _ IDateService = (*EveRefsDateService)(nil)

var earliestDate = time.Date(int(firstYearAvailable), 10, 1, 0, 0, 0, 0, time.UTC)

// lists the days in a given year, exluding anything before
// October 1, 2003, the first date with EveRefs data
func enumerateDatesInYear(year Year) []time.Time {
	if year < 2003 {
		return nil
	}
	startDate := time.Date(int(year), 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(int(year)+1, 1, 1, 0, 0, 0, 0, time.UTC)

	var dates []time.Time
	for currentDate := startDate; currentDate.Before(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		if !currentDate.Before(earliestDate) {
			dates = append(dates, currentDate)
		}
	}

	return dates
}

func filterDatesBefore(dates []time.Time, currentDate time.Time) []time.Time {
	var filteredDates []time.Time

	for _, date := range dates {
		if date.Before(currentDate.AddDate(0, 0, 1)) {
			filteredDates = append(filteredDates, date)
		}
	}

	return filteredDates
}

// generates a list of dates between 2003-10-01 and the date
// returned by the EveRefsDateservice's datefn
func (d *EveRefsDateService) EnumerateDates() []time.Time {
	currentDate := d.datefn()
	numYears := currentDate.Year() - earliestDate.Year() + 1

	var dates []time.Time

	for i := range numYears {
		currentYear := Year(earliestDate.Year() + i)
		currentYearDates := enumerateDatesInYear(currentYear)
		dates = append(dates, currentYearDates...)
	}

	return filterDatesBefore(dates, currentDate)
}

// constructor for the dateservice
// datefn should, generally, be a function which returns the current date
func NewEveRefsDateService(datefn func() time.Time) IDateService {
	return &EveRefsDateService{
		datefn: datefn,
	}
}
