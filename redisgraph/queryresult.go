package redisgraph

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type ColumnType int

const (
	COLUMN_UNKNOWN ColumnType = iota
	COLUMN_SCALAR
	COLUMN_NODE
	COLUMN_RELATION
)

type ScalarType int

const (
	VALUE_UNKNOWN ScalarType = iota
	VALUE_NULL
	VALUE_STRING
	VALUE_INTEGER
	VALUE_BOOLEAN
	VALUE_DOUBLE
	VALUE_ARRAY
	VALUE_EDGE
	VALUE_NODE
)

type QueryResult struct {
	Rows [][]*ResultCell
	Header QueryResultHeader
}

type QueryResultHeader struct {
	names []string
	types []ColumnType
}

func createQueryResult(results interface{}) (*QueryResult, error) {
	resultSet := &QueryResult{
		Header: QueryResultHeader{
			names: make([]string, 0),
			types: make([]ColumnType, 0),
		},
	}

	values, err := redis.Values(results, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println(fmt.Sprintf("Got values: %+v\n", values))

	// Check if we've encountered an error
	if err, ok := values[len(values) - 1].(redis.Error); ok {
		return nil, err
	}

	fmt.Println("No error, parsing")

	resultSet.parseValues(values)

	return nil, nil
}

func (qr *QueryResult) parseValues(values []interface{}) {
	qr.parseHeader(values[0])
	qr.parseRecords(values[1:])
}

func (qr *QueryResult) parseHeader(rawheader interface{}) {
	header, _ := redis.Values(rawheader, nil)
	// TODO handle error

	fmt.Println(fmt.Sprintf("Checking header %+v", header))
	for _, col := range header {
		fmt.Println(fmt.Sprintf("Got column %+v", col))
		c, _ := redis.Values(col, nil)
		fmt.Println(fmt.Sprintf("Got c %+v", c))
		coltype, _ := redis.Int(c[0], nil)
		name, _ := redis.String(c[1], nil)

		fmt.Println(fmt.Sprintf("Got type %+v", coltype))
		fmt.Println(fmt.Sprintf("Got column names %+v", name))

		qr.Header.types = append(qr.Header.types, ColumnType(coltype))
		qr.Header.names = append(qr.Header.names, name)
	}
}

func (qr *QueryResult) parseRecords(rawresults []interface{}) {
	records, _ := redis.Values(rawresults, nil)
	// TODO handle error

	qr.Rows = make([][]*ResultCell, len(records))

	for i, record := range records {
		cells, _ := redis.Values(record, nil)
		// TODO handle error
		record := make([]*ResultCell, len(cells))

		for idx, cell := range cells {
			// Only going to support scalar types for now
			coltype := qr.Header.types[idx]
			switch coltype {
			case COLUMN_SCALAR:
				s, _ := redis.Values(cell, nil)
				t, _ := redis.Int(s[0], nil)
				record[idx] = &ResultCell{RawValue: s[1], ColType: ScalarType(t)}
				break
			case COLUMN_NODE:
				// encountered node, print warning
				break
			case COLUMN_RELATION:
				// encountered relation, print warning
				break
			default:
				// return error - unknown column type
			}
		}

		qr.Rows[i] = record
	}
}
