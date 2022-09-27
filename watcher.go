package etcdwatcher

import (
	"context"
	"fmt"
	"log"
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
