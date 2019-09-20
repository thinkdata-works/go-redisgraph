package rgraphquery

import (
	"errors"
	"github.com/gomodule/redigo/redis"
)

type QueryResult struct {
	Rows [][]ResultCell
	Headers []string
}

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
	headers, err := redis.Values(rawheader, nil)
	if err != nil {
		return err
	}

	qr.Headers = make([]string, len(headers))

	for i, header := range headers {
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
