package etcdwatcher_test

import (
	"testing"

	"github.com/fregie/etcdwatcher"
)

func TestStringSlice(t *testing.T) {
	s := etcdwatcher.NewStringSlice("test", []string{"default"})
	if s.Value()[0] != "default" {
		t.Errorf("s.Value()[0] != \"default\"")
	}
	s.Parse([]byte(`["test1","test2"]`))
	if s.Value()[0] != "test1" {
		t.Errorf("s.Value()[0] != \"test1\"")
	}
	if s.Value()[1] != "test2" {
		t.Errorf("s.Value()[1] != \"test2\"")
	}
	s.SetDefault()
	if s.Value()[0] != "default" {
		t.Errorf("s.Value()[0] != \"default\"")
	}
}
