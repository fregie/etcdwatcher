package etcdwatcher_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/fregie/etcdwatcher"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	etcdAddr   = "http://localhost:52379"
	testPrefix = "/test"
)

func TestPrefixWatcher(t *testing.T) {
	waitEtcd(etcdAddr, time.Second*5)
	etcdCli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdAddr},
		DialTimeout: 5 * time.Second,
	})
	checkErr(t, err)
	_, err = etcdCli.Delete(context.Background(), testPrefix, clientv3.WithPrefix())
	checkErr(t, err)

	watcher, err := etcdwatcher.NewPrefixWatcher([]string{etcdAddr})
	checkErr(t, err)
	watchMap := etcdwatcher.NewMap(testPrefix, func(key string) etcdwatcher.Item {
		return etcdwatcher.NewString(key, "")
	})
	watcher.Watch(watchMap)
	_, ok := watchMap.GetItem("key1")
	if ok {
		t.Errorf("ok before set key1")
	}
	err = watchMap.SaveToEtcd(context.Background(), etcdCli, "key1", "value1")
	checkErr(t, err)
	time.Sleep(time.Second)
	i, ok := watchMap.GetItem("key1")
	if !ok {
		t.Errorf("!ok after set key1")
	}
	if i.(*etcdwatcher.String).Value() != "value1" {
		t.Errorf("i.(*etcdwatcher.String).Value != \"value1\"")
	}
	err = watchMap.DeleteFromEtcd(context.Background(), etcdCli, "key1")
	checkErr(t, err)
	time.Sleep(time.Second)
	_, ok = watchMap.GetItem("key1")
	if ok {
		t.Errorf("ok after delete key1")
	}
}

func waitEtcd(addr string, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		if ctx.Err() != nil {
			return
		}
		_, err := http.Get(addr)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
}

// checkErr is a helper function to check error
func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("err: %v", err)
	}
}
