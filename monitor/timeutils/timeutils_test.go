package timeutils_test

import (
	"testing"
	"time"

	"github.com/hoodnoah/eve_market/monitor/timeutils"
)

func sliceSame(s1 []int, s2 []int) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range len(s1) {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}

func TestEnumerateYearsToPresent(t *testing.T) {
	startYear := 2003

	expectedResult := []int{2003, 2004, 2005, 2006, 2007, 2008, 2009, 2010, 2011, 2012, 2013, 2014, 2015, 2016, 2017, 2018, 2019, 2020, 2021, 2022, 2023, 2024}

	actualResult := timeutils.EnumerateYearsToPresent(startYear)

	if !sliceSame(expectedResult, actualResult) {
		t.Fatal("Slices unequal")
	}
}

func TestEnumerateDatesInYearReturnsFirstYearCorrectly(t *testing.T) {
	year := 2003
	dates := timeutils.EnumerateDatesInYear(year)

	if len(dates) != 92 {
		t.Fatalf("Expected 92 days, got: %d", len(dates))
	}
}

func TestEnumerateDatesInYearReturns365ForNonLeapYear(t *testing.T) {
	year := 2005
	dates := timeutils.EnumerateDatesInYear(year)

	if len(dates) != 365 {
		t.Fatalf("Expected 365 days, rec'd %d", len(dates))
	}
}

func TestEnumerateDatesInYearReturns366ForLeapYear(t *testing.T) {
	year := 2004
	dates := timeutils.EnumerateDatesInYear(year)

	if len(dates) != 366 {
		t.Fatalf("Expected 366 days, rec'd %d", len(dates))
	}
}

func TestGetFileURL(t *testing.T) {
	date := time.Date(2003, 10, 1, 0, 0, 0, 0, time.UTC)

	expectedResult := "2003/market-history-2003-10-01.csv.bz2"

	result := timeutils.GetDayFileURL(date)

	if !(result == expectedResult) {
		t.Fatalf("Expected %s, received: %s", expectedResult, result)
	}
}
