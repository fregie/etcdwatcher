package etcdwatcher

import (
	"strconv"
	"sync/atomic"
)

type Bool struct {
	key          string
	defaultValue int32
	value        int32
}

func NewBool(key string, defaultBool bool) *Bool {
	b := &Bool{
		key: key,
	}
	var value int32
	if defaultBool {
		value = 1
	} else {
		value = 0
	}
	atomic.StoreInt32(&b.value, value)
	atomic.StoreInt32(&b.defaultValue, value)
	return b
}

func (b *Bool) SetValue(value bool) {
	if value {
		atomic.StoreInt32(&b.value, 1)
	} else {
		atomic.StoreInt32(&b.value, 0)
	}
}

func (b *Bool) SetDefault() {
	atomic.StoreInt32(&b.value, atomic.LoadInt32(&b.defaultValue))
}

func (b *Bool) Parse(data []byte) error {
	v, err := strconv.ParseBool(string(data))
	if err != nil {
		return err
	}
	b.SetValue(v)
	return nil
}

func (b *Bool) Key() string {
	return b.key
}

func (b *Bool) Value() bool {
	return atomic.LoadInt32(&b.value) == 1
}
