package etcdwatcher

import "strconv"

type Bool struct {
	key          string
	defaultValue bool
	value        bool
}

func NewBool(key string, defaultBool bool) *Bool {
	return &Bool{
		key:          key,
		defaultValue: defaultBool,
		value:        defaultBool,
	}
}

func (b *Bool) SetDefault() {
	b.value = b.defaultValue
}

func (b *Bool) Parse(data []byte) error {
	v, err := strconv.ParseBool(string(data))
	if err != nil {
		return err
	}
	b.value = v
	return nil
}

func (b *Bool) Key() string {
	return b.key
}

func (b *Bool) Value() bool {
	return b.value
}
