package datadateservice_test

import (
	"fmt"
	"testing"
	"time"

	dds "github.com/hoodnoah/eve_market/monitor/datadateservice"
)

// wraps a provided date in a getter function
// to match the signature req'd by the
// DataDateService struct
func getMockDateFn(date time.Time) func() time.Time {
	return func() time.Time {
		return date
	}
}

func Test_returns_correct_num_years_2003(t *testing.T) {
	dateFn := getMockDateFn(time.Date(2003, 12, 31, 0, 0, 0, 0, time.UTC))

	service := dds.NewDataDateService(dateFn)

	dataYears := service.EnumerateDataYears()

	if len(dataYears) != 1 {
		t.Fatalf("Expected a single year, received %d", len(dataYears))
	}
}

func Test_returns_correct_num_years_2024(t *testing.T) {
	dateFn := getMockDateFn(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))

	service := dds.NewDataDateService(dateFn)

	dataYears := service.EnumerateDataYears()

	if len(dataYears) != 22 {
		t.Fatalf("Expected 22 years, received %d", len(dataYears))
	}
}

func Test_returns_correct_num_days_2003(t *testing.T) {
	dateFn := getMockDateFn(time.Date(2003, 12, 31, 0, 0, 0, 0, time.UTC))

	service := dds.NewDataDateService(dateFn)

	dataYears := service.EnumerateDataYears()

	year2003 := dataYears[0]

	if len(year2003.Dates) != 92 {
		t.Fatalf("Expected 92 days, received %d", len(year2003.Dates))
	}
}

func Test_returns_correct_num_days_leap_year(t *testing.T) {
	dateFn := getMockDateFn(time.Date(2005, 12, 31, 0, 0, 0, 0, time.UTC))

	service := dds.NewDataDateService(dateFn)

	dataYears := service.EnumerateDataYears()

	year2004 := dataYears[1]

	if len(year2004.Dates) != 366 {
		t.Fatalf("Expected 366 days, received %d", len(year2004.Dates))
	}
}

func Test_returns_correct_num_days_for_partial_year(t *testing.T) {
	dateFn := getMockDateFn(time.Date(2010, 6, 15, 0, 0, 0, 0, time.UTC))

	service := dds.NewDataDateService(dateFn)

	dataYears := service.EnumerateDataYears()

	year := dataYears[len(dataYears)-1]

	for _, day := range year.Dates {
		fmt.Println(day.Url)
	}

	if len(year.Dates) != 166 {
		t.Fatalf("Expected 166 days, received %d", len(year.Dates))
	}
}

func Test_returns_nothing_for_years_too_early(t *testing.T) {
	dateFn := getMockDateFn(time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC))

	service := dds.NewDataDateService(dateFn)

	dataYears := service.EnumerateDataYears()

	if len(dataYears) != 0 {
		t.Fatalf("Expected no years, received %d", len(dataYears))
	}
}
