package etcdwatcher_test

import (
	"testing"

	"github.com/fregie/etcdwatcher"
)

func TestString(t *testing.T) {
	s := etcdwatcher.NewString("test", "default")
	if s.Value() != "default" {
		t.Errorf("s.Value() != \"default\"")
	}
	s.Parse([]byte("test"))
	if s.Value() != "test" {
		t.Errorf("s.Value() != \"test\"")
	}
	s.SetDefault()
	if s.Value() != "default" {
		t.Errorf("s.Value() != \"default\"")
	}
}
