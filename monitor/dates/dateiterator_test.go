package dates_test

import (
	"testing"
	"time"

	"github.com/hoodnoah/eve_market/monitor/dates"
)

func setup(mockCurrentDate time.Time) *dates.DateIterator {
	datefn := func() time.Time {
		return mockCurrentDate
	}

	return dates.NewDateIterator(datefn)
}

func TestDateIterator(t *testing.T) {
	t.Run("it should return 1 date if datefn returns 2003-10-2", func(t *testing.T) {
		firstDate := time.Date(2003, 10, 1, 0, 0, 0, 0, time.UTC)
		currentDate := time.Date(2003, 10, 2, 0, 0, 0, 0, time.UTC)
		dateIterator := setup(currentDate)

		datesList := make([]time.Time, 0, 1)
		for dateIterator.IsNextReady() {
			datesList = append(datesList, dateIterator.Next())
		}

		if len(datesList) != 1 {
			t.Fatalf("Expected 1 date ready, received %d", len(datesList))
		}

		if !datesList[0].Equal(firstDate) {
			t.Fatalf("Expected only %s, received %s", currentDate.String(), datesList[0].String())
		}
	})

	t.Run("it should return no dates if datefn returns any time on 2003-10-01", func(t *testing.T) {
		firstDate := time.Date(2003, 10, 1, 15, 3, 2, 1, time.UTC)
		dateIterator := setup(firstDate)

		datesList := make([]time.Time, 0, 1)
		for dateIterator.IsNextReady() {
			datesList = append(datesList, dateIterator.Next())
		}

		if len(datesList) != 0 {
			t.Fatalf("Expected 0 date ready, received %d", len(datesList))
		}
	})
}
