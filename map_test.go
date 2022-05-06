package etcdwatcher

import "testing"

func TestMap(t *testing.T) {
	m := NewMap("test", func(key string) Item {
		return NewBool(key, true)
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
	if v.(*Bool).Value() != false {
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
