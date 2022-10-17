package etcdwatcher

import (
	"sync/atomic"
	"time"
)

type Duration struct {
	key          string
	defaultValue int64
	value        int64
}

func NewDuration(key string, defaultDuration time.Duration) *Duration {
	return &Duration{
		key:          key,
		defaultValue: int64(defaultDuration),
		value:        int64(defaultDuration),
	}
}

func (p *Duration) SetDefault() {
	atomic.StoreInt64(&p.value, p.defaultValue)
}

func (d *Duration) Parse(data []byte) error {
	val, err := time.ParseDuration(string(data))
	if err != nil {
		return err
	}
	addr := &d.value
	atomic.StoreInt64(addr, int64(val))
	return nil
}

func (d *Duration) Key() string {
	return d.key
}

func (d *Duration) Value() time.Duration {
	return time.Duration(atomic.LoadInt64(&d.value))
}
