package graph

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type ResultCell struct {
	value interface{}

}

func (r *ResultCell) ToString() (string, error) {
	return redis.String(r.value, nil)
}

func (r *ResultCell) ToInt() (int, error) {
	return redis.Int(r.value, nil)
}

func (r *ResultCell) ToInt64() (int64, error) {
	return redis.Int64(r.value, nil)
}

func (r *ResultCell) ToBool() (bool, error) {
	return redis.Bool(r.value, nil)
}

func (r *ResultCell) ToFloat64() (float64, error) {
	return redis.Float64(r.value, nil)
}

func (r *ResultCell) IsNull() bool {
	switch r.value.(type) {
	case nil:
		return true
	default:
		return false
	}
}

func (r *ResultCell) ToArray() ([]ResultCell, error) {
	cells := r.value.([]interface{})
	array := make([]ResultCell, len(cells))
	for i, c := range(cells) {
		array[i] = ResultCell{value: c}
	}
	return array, nil
}

func (r *ResultCell) Inspect() string {
	s, _ := r.ToString()
	i, _ := r.ToInt()
	i64, _ := r.ToInt64()
	b, _ := r.ToBool()
	return fmt.Sprintf("{ value: %+v, String: %s, Int: %d, Int64: %d, Bool: %+v }", r.value, s, i, i64, b)
}