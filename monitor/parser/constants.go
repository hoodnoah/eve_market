package parser

var fieldNames = [8]string{
	"date",
	"region_id",
	"type_id",
	"average",
	"highest",
	"lowest",
	"volume",
	"order_count",
}

const urlTemplate = "https://data.everef.net/market-history/%04d/market-history-%04d-%02d-%02d.csv.bz2"
