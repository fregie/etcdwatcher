package etcdwatcher_test

import (
	"testing"

	"github.com/fregie/etcdwatcher"
)

func TestMap(t *testing.T) {
	m := etcdwatcher.NewMap("test", func(key string) etcdwatcher.Item {
		return etcdwatcher.NewBool(key, true)
	})
	if m.PrefixKey() != "test" {
		t.Errorf("m.PrefixKey() != \"test\"")
	}
	if err := m.ParseWithKey("key1", []byte("false")); err != nil {
		t.Errorf("err != nil")
	}
	v, ok := m.GetItem("key1")
	if !ok {
		t.Errorf("!ok")
	}
	if v.(*etcdwatcher.Bool).Value() != false {
		t.Errorf("v.Value() != false")
	}
	v, ok = m.GetItem("key2")
	if ok {
		t.Errorf("ok")
	}
	if v != nil {
		t.Errorf("v != nil")
	}
	m.DeleteKey("key1")
	v, ok = m.GetItem("key1")
	if ok {
		t.Errorf("ok")
	}
	if v != nil {
		t.Errorf("v != nil")
	}
}
