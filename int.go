package etcdwatcher

import (
	"strconv"
	"sync/atomic"
)

type Int32 struct {
	key          string
	defaultValue int32
	value        int32
}

func NewInt32(key string, defaultValue int32) *Int32 {
	return &Int32{
		key:          key,
		defaultValue: defaultValue,
		value:        defaultValue,
	}
}

func (p *Int32) SetDefault() {
	atomic.StoreInt32(&p.value, p.defaultValue)
}

func (p *Int32) Parse(data []byte) error {
	val, err := strconv.ParseInt(string(data), 10, 32)
	if err != nil {
		return err
	}
	atomic.StoreInt32(&p.value, int32(val))
	return nil
}

func (p *Int32) Key() string {
	return p.key
}

func (p *Int32) Value() int32 {
	return atomic.LoadInt32(&p.value)
}

type Int64 struct {
	key          string
	defaultValue int64
	value        int64
}

func NewInt64(key string, defaultValue int64) *Int64 {
	return &Int64{
		key:          key,
		defaultValue: defaultValue,
		value:        defaultValue,
	}
}

func (p *Int64) SetDefault() {
	atomic.StoreInt64(&p.value, p.defaultValue)
}

func (p *Int64) Parse(data []byte) error {
	val, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	atomic.StoreInt64(&p.value, val)
	return nil
}

func (p *Int64) Key() string {
	return p.key
}

func (p *Int64) Value() int64 {
	return atomic.LoadInt64(&p.value)
}
