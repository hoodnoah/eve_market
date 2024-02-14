package main

import (
	"fmt"
	"log"

	mds "github.com/hoodnoah/eve_market/monitor/marketdataservice"
)

func main() {
	const url = "https://data.everef.net/market-history/2003/market-history-2003-10-01.csv.bz2"

	downloadService := mds.NewMarketDataService(mds.DownloadFile, mds.DecompressFile, mds.ParseFile)

	result, err := downloadService.FetchAndParseCSV(url)

	if err != nil {
		log.Fatalf("Failed to download file at %s: %s", url, err)
	}

	randomRecord := result[1234]

	fmt.Printf("%+v\n", randomRecord)
}
