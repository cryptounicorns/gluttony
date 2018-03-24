package clickhouse

const (
	dbCreate = `
CREATE TABLE IF NOT EXISTS datapoints (
	country_code FixedString(2),
	os_id        UInt8,
	browser_id   UInt8,
	categories   Array(Int16),
	action_day   Date,
	action_time  DateTime
) engine=Memory
`

	dbInsert = `
INSERT INTO datapoints (country_code, os_id, browser_id, categories, action_day, action_time)
VALUES (?, ?, ?, ?, ?, ?)
`
)
