package redisgraph

import "github.com/gomodule/redigo/redis"

type ResultCell struct {
	Value interface{}
}

func (r *ResultCell) ToString() (string, error) {
	return redis.String(r.Value, nil)
}

func (r *ResultCell) ToInt() (int, error) {
	return redis.Int(r.Value, nil)
}

func (r *ResultCell) ToInt64() (int64, error) {
	return redis.Int64(r.Value, nil)
}

func (r *ResultCell) ToBool() (bool, error) {
	return redis.Bool(r.Value, nil)
}

func (r *ResultCell) ToFloat64() (float64, error) {
	return redis.Float64(r.Value, nil)
}