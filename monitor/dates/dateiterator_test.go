package dates_test

import (
	"testing"
	"time"

	"github.com/hoodnoah/eve_market/monitor/dates"
)

type TestSetup struct {
	Iterator      *dates.DateIterator
	OutputChannel chan time.Time
}

func setup(mockCurrentDate time.Time) *TestSetup {
	datefn := func() time.Time {
		return mockCurrentDate
	}
	channel := make(chan time.Time, 1)

	iterator := dates.NewDateIterator(datefn, channel)

	return &TestSetup{
		Iterator:      iterator,
		OutputChannel: channel,
	}
}

func dateChanToList(channel chan time.Time) []time.Time {
	var dates []time.Time
	timeout := time.After(250 * time.Millisecond)

	for {
		select {
		case d := <-channel:
			dates = append(dates, d)
		case <-timeout:
			return dates
		}
	}
}

func TestDateIterator(t *testing.T) {
	t.Run("it should return 1 date if datefn returns 2003-10-2", func(t *testing.T) {
		firstDate := time.Date(2003, 10, 1, 0, 0, 0, 0, time.UTC)
		currentDate := time.Date(2003, 10, 2, 0, 0, 0, 0, time.UTC)

		// setup test dependencies
		testSetup := setup(currentDate)
		testSetup.Iterator.Start()
		outputChannel := testSetup.OutputChannel
		defer close(outputChannel)

		// extract dates into a list for easy comparison
		datesList := dateChanToList(outputChannel)

		if len(datesList) != 1 {
			t.Fatalf("Expected 1 date ready, received %d", len(datesList))
		}

		if !datesList[0].Equal(firstDate) {
			t.Fatalf("Expected only %s, received %s", currentDate.String(), datesList[0].String())
		}
	})

	t.Run("it should return no dates if datefn returns any time on 2003-10-01", func(t *testing.T) {
		firstDate := time.Date(2003, 10, 1, 15, 3, 2, 1, time.UTC)
		testSetup := setup(firstDate)
		testSetup.Iterator.Start()
		defer close(testSetup.OutputChannel)

		datesList := dateChanToList(testSetup.OutputChannel)

		if len(datesList) != 0 {
			t.Fatalf("Expected 0 date ready, received %d", len(datesList))
		}
	})

	t.Run("it should produce 92 dates, when nowFn says it's 2004-01-01", func(t *testing.T) {
		firstDate := time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC)
		testSetup := setup(firstDate)
		testSetup.Iterator.Start()
		defer close(testSetup.OutputChannel)

		datesList := dateChanToList(testSetup.OutputChannel)

		// should contain the right number of days
		if len(datesList) != 92 {
			t.Fatalf("Expected 92 days, received %d", len(datesList))
		}

		// first date should be 2003-10-01
		if !datesList[0].Equal(time.Date(2003, 10, 1, 0, 0, 0, 0, time.UTC)) {
			t.Fatalf("Expected first date to be 2003-10-01, received %s", datesList[0].String())
		}

		// last date should be 2003-12-31
		if !datesList[len(datesList)-1].Equal(time.Date(2003, 12, 31, 0, 0, 0, 0, time.UTC)) {
			t.Fatalf("Expected last date to be 2003-12-31, received %s", datesList[len(datesList)-1].String())
		}
	})

	t.Run("it should produce no dates if the date function returns a date before 2003-10-01", func(t *testing.T) {
		firstDate := time.Date(2003, 1, 12, 0, 0, 0, 0, time.UTC)
		testSetup := *setup(firstDate)
		testSetup.Iterator.Start()
		defer close(testSetup.OutputChannel)

		datesList := dateChanToList(testSetup.OutputChannel)

		if len(datesList) != 0 {
			t.Fatalf("Expected 0 dates returned, received %d", len(datesList))
		}
	})

	t.Run("it should produce a large number of dates in a performant manner, e.g. from 2003-10-01 through 2024-12-31", func(t *testing.T) {
		firstDate := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

		testSetup := setup(firstDate)
		testSetup.Iterator.Start()
		defer close(testSetup.OutputChannel)

		datesList := dateChanToList(testSetup.OutputChannel)

		if len(datesList) != 7762 {
			t.Fatalf("Expected 7762 dates, received %d", len(datesList))
		}
	})
}
