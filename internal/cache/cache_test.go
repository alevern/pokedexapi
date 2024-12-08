package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestGetNotFound(t *testing.T) {
	const interval = 5 * time.Second
	cache := NewCache(interval)
	data := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}
	for _, c := range data {
		cache.Add(c.key, c.val)
	}
	existing := data[0].key
	unknown := "https://example.com/unknown"
	_, ok := cache.Get(unknown)
	if ok {
		t.Errorf("Did not expected to find key %s", unknown)
		return
	}
	_, ok = cache.Get(existing)
	if !ok {
		t.Errorf("expected to find key %s", existing)
		return
	}
}

func TestOverrideAdd(t *testing.T) {
	const interval = 5 * time.Second
	cache := NewCache(interval)
	data := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}
	for _, c := range data {
		cache.Add(c.key, c.val)
	}
	existing := data[0].key
	newData := []byte("moretestdata")
	cache.Add(existing, newData)
	val, ok := cache.Get(existing)
	if !ok {
		t.Errorf("expected to find key %s", existing)
		return
	}
	if string(val) != string(newData) {
		t.Errorf("expected to find value %s", string(val))
		return
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}
