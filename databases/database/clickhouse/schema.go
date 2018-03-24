package clickhouse

const (
	dbCreate = `
CREATE TABLE IF NOT EXISTS datapoints (
	timestamp UInt64,
	open      UInt64,
	high      UInt64,
	low       UInt64,
	close     UInt64,
	volume    UInt64,
) engine=Memory
`

	dbInsert = `
INSERT INTO datapoints (timestamp, open, high, low, close, volume)
VALUES (?, ?, ?, ?, ?, ?)
`
)
