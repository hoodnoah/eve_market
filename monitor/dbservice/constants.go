package dbservice

// CONSTANTS
var marketDataTableTemplate = `
	CREATE TABLE IF NOT EXISTS market_data (
		date_id INTEGER NOT NULL,
		region_id INTEGER UNSIGNED NOT NULL,
		type_id INTEGER UNSIGNED NOT NULL,
		average DECIMAL(20, 2) NOT NULL,
		highest DECIMAL(20, 2) NOT NULL,
		lowest DECIMAL(20, 2) NOT NULL,
		volume INTEGER UNSIGNED NOT NULL,
		order_count INTEGER UNSIGNED NOT NULL,

		PRIMARY KEY(date_id, region_id, type_id),
		FOREIGN KEY (date_id)
			REFERENCES completed_dates(id)
			ON UPDATE CASCADE
			ON DELETE CASCADE
	);
`

var completedDatesTableTemplate = `
	CREATE TABLE IF NOT EXISTS completed_dates (
		id INT NOT NULL AUTO_INCREMENT,
		date DATE,

		PRIMARY KEY (id)
	);
`

var insertCompletedDateTemplate = `
	INSERT INTO completed_dates
		(date)
	VALUES
		(?)
	ON DUPLICATE KEY UPDATE date=date
`

var insertManyTemplate = `
	INSERT INTO market_data
		(date_id, region_id, type_id, average, highest, lowest, volume, order_count)
	VALUES
		%s
	ON DUPLICATE KEY UPDATE date_id=date_id;
`

const MAXCHUNKSIZE = 2000
