package etcdwatcher

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Watcher struct {
	etcdCli *clientv3.Client
	logger  *log.Logger
}

func NewWatcher(endpoints []string) (*Watcher, error) {
	etcdCli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &Watcher{
		etcdCli: etcdCli,
		logger:  log.Default(),
	}, nil
}

func (w *Watcher) WatchItems(items []Item) error {
	for _, item := range items {
		err := w.Watch(item)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Watcher) Watch(item Item) error {
	if item == nil {
		return fmt.Errorf("parser is nil")
	}
	go func() {
		for {
			var rev int64 = 0
			rsp, err := w.etcdCli.Get(context.Background(), item.Key())
			if err != nil {
				w.logger.Printf("etcd get error: %v", err)
			} else {
				if len(rsp.Kvs) > 0 {
					// w.logger.Printf("%s: %s", item.Key(), string(rsp.Kvs[0].Value))
					err := item.Parse(rsp.Kvs[0].Value)
					if err != nil {
						w.logger.Printf("parser error: %v", err)
					}
				}
				rev = rsp.Header.Revision + 1
			}
			watcher := clientv3.NewWatcher(w.etcdCli)
			rspChan := watcher.Watch(context.Background(), item.Key(), clientv3.WithRev(rev))
			for rsp := range rspChan {
				for _, ev := range rsp.Events {
					switch ev.Type {
					case clientv3.EventTypePut:
						err := item.Parse(ev.Kv.Value)
						if err != nil {
							w.logger.Printf("parser error: %v", err)
						}
					case clientv3.EventTypeDelete:
						item.SetDefault()
					}
				}
			}
			time.Sleep(3 * time.Second)
		}
	}()
	return nil
}

type Item interface {
	// Return the key watching
	Key() string
	// Parse the data when the value is changed
	Parse([]byte) error
	// Set the default value
	SetDefault()
}

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
