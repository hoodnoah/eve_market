package dateservice_test

import (
	"testing"
	"time"

	dds "github.com/hoodnoah/eve_market/monitor/dateservice"
)

func filterDates(dates []time.Time, comparator func(time.Time) bool) []time.Time {
	filteredDates := make([]time.Time, 0)
	for _, date := range dates {
		if comparator(date) {
			filteredDates = append(filteredDates, date)
		}
	}
	return filteredDates
}

func countYears(dates []time.Time) int {
	isIn := func(year int, years []int) bool {
		for _, y := range years {
			if year == y {
				return true
			}
		}
		return false
	}

	distinctYears := make([]int, 0)
	for _, date := range dates {
		if !isIn(date.Year(), distinctYears) {
			distinctYears = append(distinctYears, date.Year())
		}
	}

	return len(distinctYears)
}

// wraps a provided date in a getter function
// to match the signature req'd by the
// DataDateService struct
func getMockDateFn(date time.Time) func() time.Time {
	return func() time.Time {
		return date
	}
}

func setup(date time.Time) struct {
	Service dds.IDateService
} {
	datefn := func() time.Time {
		return date
	}

	svc := dds.NewDatesToPresentService(datefn)

	return struct {
		Service dds.IDateService
	}{
		Service: svc,
	}
}

func Test_returns_correct_num_years_2003(t *testing.T) {
	// arrange
	fixture := setup(time.Date(2003, 12, 31, 0, 0, 0, 0, time.UTC))
	service := fixture.Service

	// act
	dates := service.EnumerateDates()
	numYears := countYears(dates)

	// assert
	if numYears != 1 {
		t.Fatalf("Expected a single year, received %d", numYears)
	}
}

func Test_returns_correct_num_years_2024(t *testing.T) {
	// arrange
	fixture := setup(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	service := fixture.Service

	// act
	dates := service.EnumerateDates()
	numYears := countYears(dates)

	// assert
	if numYears != 22 {
		t.Fatalf("Expected 22 years, received %d", numYears)
	}
}

func Test_returns_correct_num_days_2003(t *testing.T) {
	// arrange
	fixture := setup(time.Date(2003, 12, 31, 0, 0, 0, 0, time.UTC))
	service := fixture.Service

	comparator := func(date time.Time) bool {
		return date.Year() == 2003
	}

	// act
	dates := service.EnumerateDates()
	year2003 := filterDates(dates, comparator)

	// assert
	if len(year2003) != 92 {
		t.Fatalf("Expected 92 days, received %d", len(year2003))
	}
}

func Test_returns_correct_num_days_leap_year(t *testing.T) {
	fixture := setup(time.Date(2005, 12, 31, 0, 0, 0, 0, time.UTC))
	service := fixture.Service
	comparator := func(date time.Time) bool {
		return date.Year() == 2004
	}

	dates := service.EnumerateDates()
	year2004 := filterDates(dates, comparator)

	if len(year2004) != 366 {
		t.Fatalf("Expected 366 days, received %d", len(year2004))
	}
}

func Test_returns_correct_num_days_for_partial_year(t *testing.T) {
	fixture := setup(time.Date(2010, 6, 15, 0, 0, 0, 0, time.UTC))
	service := fixture.Service
	comparator := func(date time.Time) bool {
		return date.Year() == 2010
	}

	dates := filterDates(service.EnumerateDates(), comparator)

	if len(dates) != 166 {
		t.Fatalf("Expected 166 days, received %d", len(dates))
	}
}

func Test_returns_nothing_for_years_too_early(t *testing.T) {
	fixture := setup(time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC))
	service := fixture.Service

	dates := service.EnumerateDates()

	if len(dates) != 0 {
		t.Fatalf("Expected no years, received %d", len(dates))
	}
}
