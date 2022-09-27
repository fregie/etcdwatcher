package etcdwatcher_test

import (
	"context"
	"testing"
	"time"

	"github.com/fregie/etcdwatcher"
	"github.com/stretchr/testify/suite"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type WatchTestSuite struct {
	suite.Suite
}

func TestWatcherSuite(t *testing.T) {
	suite.Run(t, new(WatchTestSuite))
}

func (s *WatchTestSuite) TestWatcher() {
	waitEtcd(etcdAddr, time.Second*5)
	etcdCli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdAddr},
		DialTimeout: 5 * time.Second,
	})
	s.Nil(err)
	_, err = etcdCli.Delete(context.Background(), "/testKey", clientv3.WithPrefix())
	s.Nil(err)
	watcher, err := etcdwatcher.NewWatcher([]string{etcdAddr})
	s.Nil(err)
	int32Value := etcdwatcher.NewInt32("/testKey/int32", 10)
	int64Value := etcdwatcher.NewInt64("/testKey/int64", 10)
	stringValue := etcdwatcher.NewString("/testKey/string", "default")
	durationValue := etcdwatcher.NewDuration("/testKey/duration", time.Second)
	boolValue := etcdwatcher.NewBool("/testKey/bool", false)
	err = watcher.WatchItems([]etcdwatcher.Item{
		int32Value,
		int64Value,
		stringValue,
		durationValue,
		boolValue,
	})
	s.Nil(err)
	s.Equal(int32(10), int32Value.Value())
	s.Equal(int64(10), int64Value.Value())
	s.Equal("default", stringValue.Value())
	s.Equal(time.Second, durationValue.Value())
	s.Equal(false, boolValue.Value())
	etcdCli.Put(context.Background(), "/testKey/int32", "20")
	etcdCli.Put(context.Background(), "/testKey/int64", "20")
	etcdCli.Put(context.Background(), "/testKey/string", "string")
	etcdCli.Put(context.Background(), "/testKey/duration", "2s")
	etcdCli.Put(context.Background(), "/testKey/bool", "true")
	time.Sleep(300 * time.Millisecond)
	s.Equal(int32(20), int32Value.Value())
	s.Equal(int64(20), int64Value.Value())
	s.Equal("string", stringValue.Value())
	s.Equal(time.Second*2, durationValue.Value())
	s.Equal(true, boolValue.Value())
}
