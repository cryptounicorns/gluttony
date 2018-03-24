package record

type Record = interface{}

type DataRecord struct {
	Timestamp uint64
	Open      uint64
	High      uint64
	Low       uint64
	Close     uint64
	Volume    uint64
}
