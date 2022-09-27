package etcdwatcher_test

import (
	"testing"
	"time"

	"github.com/fregie/etcdwatcher"
)

func TestDuration(t *testing.T) {
	d := etcdwatcher.NewDuration("test", 10*time.Second)
	if d.Value() != 10*time.Second {
		t.Errorf("d.Value() != 10s")
	}
	d.Parse([]byte("20s"))
	if d.Value() != 20*time.Second {
		t.Errorf("d.Value() != 20s")
	}
	d.SetDefault()
	if d.Value() != 10*time.Second {
		t.Errorf("d.Value() != 10")
	}
	err := d.Parse([]byte("test"))
	if err == nil {
		t.Errorf("err == nil")
	}
	if d.Value() != 10*time.Second {
		t.Errorf("d.Value() != 10")
	}
}
