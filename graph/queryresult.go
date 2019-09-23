package graph

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type QueryResult struct {
	Rows [][]ResultCell
	Headers []string
	HeaderTypes []ColumnTypes
}

type ColumnTypes int

const (
	COLUMN_UNKNOWN ColumnTypes = iota
	COLUMN_SCALAR
	COLUMN_NODE
	COLUMN_RELATION
)

type ScalarTypes int

const (
	VALUE_UNKNOWN ScalarTypes = iota
	VALUE_NULL
	VALUE_STRING
	VALUE_INTEGER
	VALUE_BOOLEAN
	VALUE_DOUBLE
	VALUE_ARRAY
	VALUE_EDGE
	VALUE_NODE
)

func createQueryResult(results interface{}) (*QueryResult, error) {
	resultSet := &QueryResult{}

	values, err := redis.Values(results, nil)
	if err != nil {
		return nil, err
	}

	// Check if we've encountered an error
	if err, ok := values[len(values) - 1].(redis.Error); ok {
		return nil, err
	}

	// First row will always be headers
	if len(values) > 1 {
		err := resultSet.parseValues(values)
		if err != nil {
			return nil, err
		}
	}

	return resultSet, nil
}

func (qr *QueryResult) parseValues(values []interface{}) error {
	err := qr.parseHeader(values[0])
	if err != nil {
		return err
	}

	return qr.parseRecords(values[1])
}

func (qr *QueryResult) parseHeader(rawheader interface{}) error {
	fmt.Println(fmt.Sprintf("Parsing raw header: %+v", rawheader))
	headers, err := redis.Values(rawheader, nil)
	if err != nil {
		return err
	}

	qr.Headers = make([]string, len(headers))

	for i, header := range headers {
		fmt.Println(fmt.Sprintf("Stringifying header: %+v", header))
		name, err := redis.String(header, nil)
		if err != nil {
			return err
		}

		qr.Headers[i] = name
	}

	return nil
}

func (qr *QueryResult) parseRecords(rawresults interface{}) error {
	records, err := redis.Values(rawresults, nil)
	if err != nil {
		return err
	}

	qr.Rows = make([][]ResultCell, len(records))

	for i, record := range records {

		qr.Rows[i] = make([]ResultCell, len(record.([]interface{})))

		for j, cell := range record.([]interface{}) {
			qr.Rows[i][j] = ResultCell{value: cell}
		}
	}

	return nil
}

// Returns a set of pointers to the original cells
func (qr *QueryResult) CellsFor(header string) ([]*ResultCell, error) {
	i := -1

	// get index of header
	for j, headerv := range(qr.Headers) {
		if headerv == header {
			i = j
		}
	}

	if i < 0 {
		return nil, errors.New("Header not found")
	}

	cells := make([]*ResultCell, len(qr.Rows))
	for j, row := range(qr.Rows) {
		cells[j] = &row[i]
	}

	return cells, nil
}
