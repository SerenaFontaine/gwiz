package gwiz

import (
	"slices"
	"sort"
	"testing"
)

func TestMapState_SetAndGet(t *testing.T) {
	s := newState()
	s.Set("name", "test")
	v, ok := s.Get("name")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if v != "test" {
		t.Fatalf("expected 'test', got %v", v)
	}
}

func TestMapState_GetMissing(t *testing.T) {
	s := newState()
	_, ok := s.Get("missing")
	if ok {
		t.Fatal("expected key to not exist")
	}
}

func TestMapState_GetString(t *testing.T) {
	s := newState()
	s.Set("name", "hello")
	if s.GetString("name") != "hello" {
		t.Fatalf("expected 'hello', got %q", s.GetString("name"))
	}
	if s.GetString("missing") != "" {
		t.Fatal("expected empty string for missing key")
	}
	s.Set("num", 42)
	if s.GetString("num") != "" {
		t.Fatal("expected empty string for wrong type")
	}
}

func TestMapState_GetBool(t *testing.T) {
	s := newState()
	s.Set("flag", true)
	if !s.GetBool("flag") {
		t.Fatal("expected true")
	}
	if s.GetBool("missing") {
		t.Fatal("expected false for missing key")
	}
}

func TestMapState_GetInt(t *testing.T) {
	s := newState()
	s.Set("count", 42)
	if s.GetInt("count") != 42 {
		t.Fatalf("expected 42, got %d", s.GetInt("count"))
	}
	if s.GetInt("missing") != 0 {
		t.Fatal("expected 0 for missing key")
	}
}

func TestMapState_GetStringSlice(t *testing.T) {
	s := newState()
	s.Set("items", []string{"a", "b", "c"})
	got := s.GetStringSlice("items")
	if !slices.Equal(got, []string{"a", "b", "c"}) {
		t.Fatalf("expected [a b c], got %v", got)
	}
	if s.GetStringSlice("missing") != nil {
		t.Fatal("expected nil for missing key")
	}
}

func TestMapState_Keys(t *testing.T) {
	s := newState()
	s.Set("b", 1)
	s.Set("a", 2)
	keys := s.Keys()
	sort.Strings(keys)
	if !slices.Equal(keys, []string{"a", "b"}) {
		t.Fatalf("expected [a b], got %v", keys)
	}
}

func TestMapState_Overwrite(t *testing.T) {
	s := newState()
	s.Set("key", "first")
	s.Set("key", "second")
	if s.GetString("key") != "second" {
		t.Fatalf("expected 'second', got %q", s.GetString("key"))
	}
}
