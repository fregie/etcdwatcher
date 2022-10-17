package etcdwatcher

import "sync"

type String struct {
	rwMutex      sync.RWMutex
	defaultValue string
	key          string
	value        string
}

func NewString(key string, defaultValue string) *String {
	return &String{
		key:          key,
		defaultValue: defaultValue,
		value:        defaultValue,
	}
}

func (p *String) SetDefault() {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	p.value = p.defaultValue
}

func (p *String) Parse(data []byte) error {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	p.value = string(data)
	return nil
}

func (p *String) Key() string {
	return p.key
}

func (p *String) Value() string {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()
	return p.value
}
