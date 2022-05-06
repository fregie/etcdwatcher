package etcdwatcher

import (
	"context"
	"fmt"
	"sync"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Map struct {
	prefixKey string
	items     sync.Map
	new       NewItem
}

type NewItem func(key string) Item

// NewMap returns a new Map.
func NewMap(prefixKey string, new NewItem) *Map {
	return &Map{
		prefixKey: prefixKey,
		new:       new,
	}
}

// PrefixKey returns the prefix key of the Map.
func (m *Map) PrefixKey() string {
	return m.prefixKey
}

// ParseWithKey parses the data with the given key.
func (m *Map) ParseWithKey(key string, data []byte) error {
	if m.new == nil {
		return fmt.Errorf("new item is nil")
	}
	item := m.new(key)
	if err := item.Parse(data); err != nil {
		return err
	}
	m.items.Store(key, item)
	return nil
}

// DeleteKey deletes the item with the given key.
func (m *Map) DeleteKey(key string) error {
	m.items.Delete(key)
	return nil
}

// GetItem returns the item with the given key.
func (m *Map) GetItem(key string) (Item, bool) {
	item, ok := m.items.Load(key)
	if !ok {
		return nil, false
	}
	return item.(Item), true
}

func (m *Map) GenWholeKey(key string) string {
	return fmt.Sprintf("%s/%s", m.prefixKey, key)
}

func (m *Map) SaveToEtcd(ctx context.Context, cli *clientv3.Client, key, data string) error {
	_, err := cli.Put(ctx, m.GenWholeKey(key), string(data))
	return err
}

func (m *Map) DeleteFromEtcd(ctx context.Context, cli *clientv3.Client, key string) error {
	_, err := cli.Delete(ctx, m.GenWholeKey(key))
	return err
}
