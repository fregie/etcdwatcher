package etcdwatcher_test

import (
	"testing"

	"github.com/fregie/etcdwatcher"
)

func TestInt32(t *testing.T) {
	i := etcdwatcher.NewInt32("test", 10)
	if i.Value() != 10 {
		t.Errorf("i.Value() != 10")
	}
	i.Parse([]byte("20"))
	if i.Value() != 20 {
		t.Errorf("i.Value() != 20")
	}
	i.SetDefault()
	if i.Value() != 10 {
		t.Errorf("i.Value() != 10")
	}
	err := i.Parse([]byte("test"))
	if err == nil {
		t.Errorf("err == nil")
	}
}

func TestInt64(t *testing.T) {
	i := etcdwatcher.NewInt64("test", 10)
	if i.Value() != 10 {
		t.Errorf("i.Value() != 10")
	}
	i.Parse([]byte("20"))
	if i.Value() != 20 {
		t.Errorf("i.Value() != 20")
	}
	i.SetDefault()
	if i.Value() != 10 {
		t.Errorf("i.Value() != 10")
	}
	err := i.Parse([]byte("test"))
	if err == nil {
		t.Errorf("err == nil")
	}
}
