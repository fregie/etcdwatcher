package etcdwatcher

import (
	"encoding/json"
	"sync"
)

type StringSlice struct {
	rwMutex      sync.RWMutex
	defaultValue []string
	key          string
	value        []string
}

func NewStringSlice(key string, defaultValue []string) *StringSlice {
	return &StringSlice{
		key:          key,
		defaultValue: defaultValue,
		value:        defaultValue,
	}
}

func (p *StringSlice) SetDefault() {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	p.value = p.defaultValue
}

func (p *StringSlice) Parse(data []byte) error {
	new := make([]string, 0)
	err := json.Unmarshal(data, &new)
	if err != nil {
		return err
	}
	p.rwMutex.Lock()
	p.value = new
	p.rwMutex.Unlock()
	return nil
}

func (p *StringSlice) Key() string {
	return p.key
}

func (p *StringSlice) Value() []string {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()
	copied := make([]string, len(p.value))
	copy(copied, p.value)
	return copied
}
