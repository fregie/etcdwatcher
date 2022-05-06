package etcdwatcher

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type PrefixWatcher struct {
	etcdCli *clientv3.Client
	logger  *log.Logger
}

func NewPrefixWatcher(endpoints []string) (*PrefixWatcher, error) {
	etcdCli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &PrefixWatcher{
		etcdCli: etcdCli,
		logger:  log.Default(),
	}, nil
}

type PrefixItem interface {
	PrefixKey() string
	ParseWithKey(key string, data []byte) error
	DeleteKey(key string) error
}

func (w *PrefixWatcher) Watch(item PrefixItem) error {
	if item == nil {
		return fmt.Errorf("parser is nil")
	}
	go func() {
		for {
			var rev int64 = 0
			rsp, err := w.etcdCli.Get(context.Background(), item.PrefixKey(), clientv3.WithPrefix())
			if err != nil {
				w.logger.Printf("etcd get error: %v", err)
			} else {
				for _, kv := range rsp.Kvs {
					subKey := strings.TrimPrefix(string(kv.Key), item.PrefixKey()+"/")
					err := item.ParseWithKey(subKey, kv.Value)
					if err != nil {
						w.logger.Printf("parse error: %v", err)
					}
				}
				rev = rsp.Header.Revision + 1
			}
			watcher := clientv3.NewWatcher(w.etcdCli)
			rspChan := watcher.Watch(context.Background(), item.PrefixKey(), clientv3.WithPrefix(), clientv3.WithRev(rev))
			for rsp := range rspChan {
				for _, ev := range rsp.Events {
					subKey := strings.TrimPrefix(string(ev.Kv.Key), item.PrefixKey()+"/")
					switch ev.Type {
					case clientv3.EventTypePut:
						err := item.ParseWithKey(subKey, ev.Kv.Value)
						if err != nil {
							w.logger.Printf("parse error: %v", err)
						}
					case clientv3.EventTypeDelete:
						err := item.DeleteKey(subKey)
						if err != nil {
							w.logger.Printf("delete error: %v", err)
						}
					}
				}
			}
			time.Sleep(3 * time.Second)
		}
	}()
	return nil
}
