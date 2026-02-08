package pokecache

import (
	"testing"
	"time"
)

func TestCreateCache(t *testing.T) {
	cache := NewCache(5 * time.Second)
	if cache.cache == nil {
		t.Error("cache is nil")
	}
}

func TestAddGet(t *testing.T) {
	cache := NewCache(5 * time.Second)

	key := "test-key"
	val := []byte("test-value")

	cache.Add(key, val)

	result, ok := cache.Get(key)
	if !ok {
		t.Error("expected to find key")
	}

	if string(result) != string(val) {
		t.Errorf("expected %s, got %s", val, result)
	}
}

func TestGetNotFound(t *testing.T) {
	cache := NewCache(5 * time.Second)

	_, ok := cache.Get("nonexistent")
	if ok {
		t.Error("expected key to not be found")
	}
}

func TestReap(t *testing.T) {
	interval := 10 * time.Millisecond
	cache := NewCache(interval)

	key := "test-key"
	val := []byte("test-value")

	cache.Add(key, val)

	// Verify it exists
	_, ok := cache.Get(key)
	if !ok {
		t.Error("expected to find key")
	}

	// Wait for it to be reaped
	time.Sleep(interval + 5*time.Millisecond)

	// Verify it's been removed
	_, ok = cache.Get(key)
	if ok {
		t.Error("expected key to be reaped")
	}
}
