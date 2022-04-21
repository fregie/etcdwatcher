package etcdwatcher

import "testing"

func TestBool(t *testing.T) {
	b := NewBool("test", true)
	if b.Value() != true {
		t.Errorf("b.Value() != true")
	}
	b.Parse([]byte("false"))
	if b.Value() != false {
		t.Errorf("b.Value() != false")
	}
	b.SetDefault()
	if b.Value() != true {
		t.Errorf("b.Value() != true")
	}
	err := b.Parse([]byte("test"))
	if err == nil {
		t.Errorf("err == nil")
	}
}
