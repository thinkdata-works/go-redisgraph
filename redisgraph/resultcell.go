package redisgraph

import "github.com/gomodule/redigo/redis"

// Wrapper for scalar types
type ResultCell struct {
	RawValue interface{}
	ColType ScalarType
}

func (r *ResultCell) ToString() string {
	s, _ := redis.String(r.RawValue, nil)
	return s
}

func (r *ResultCell) ToInt() int {
	s, _ := redis.Int(r.RawValue, nil)
	return s
}

func (r *ResultCell) ToInt64() int64 {
	s, _ := redis.Int64(r.RawValue, nil)
	return s
}

func (r *ResultCell) ToBool() bool {
	s, _ := redis.Bool(r.RawValue, nil)
	return s
}

func (r *ResultCell) ToFloat64() float64 {
	s, _ := redis.Float64(r.RawValue, nil)
	return s
}

func (r *ResultCell) IsNull() bool {
	return r.ColType == VALUE_NULL
}
