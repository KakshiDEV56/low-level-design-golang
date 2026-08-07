package lfucache

import "testing"

func TestLFUCache_Example1(t *testing.T) {
	c := NewLFUCache(2)
	c.Put(1, 1)
	c.Put(2, 2)
	if got := c.Get(1); got != 1 {
		t.Fatalf("Get(1) = %d, want 1", got)
	}
	c.Put(3, 3) // evicts key 2 (freq=1 is the LFU)
	if got := c.Get(2); got != -1 {
		t.Fatalf("Get(2) = %d, want -1", got)
	}
	if got := c.Get(3); got != 3 {
		t.Fatalf("Get(3) = %d, want 3", got)
	}
	c.Put(4, 4) // 1 and 3 tie at freq=2; 1 is LRU among them, evicts 1
	if got := c.Get(1); got != -1 {
		t.Fatalf("Get(1) = %d, want -1", got)
	}
	if got := c.Get(3); got != 3 {
		t.Fatalf("Get(3) = %d, want 3", got)
	}
	if got := c.Get(4); got != 4 {
		t.Fatalf("Get(4) = %d, want 4", got)
	}
}

func TestLFUCache_ZeroCapacity(t *testing.T) {
	c := NewLFUCache(0)
	c.Put(1, 1)
	if got := c.Get(1); got != -1 {
		t.Fatalf("Get(1) with capacity 0 = %d, want -1", got)
	}
}

func TestLFUCache_CapacityOne(t *testing.T) {
	c := NewLFUCache(1)
	c.Put(1, 1)
	c.Put(2, 2) // evicts key 1 immediately
	if got := c.Get(1); got != -1 {
		t.Fatalf("Get(1) = %d, want -1", got)
	}
	if got := c.Get(2); got != 2 {
		t.Fatalf("Get(2) = %d, want 2", got)
	}
}

func TestLFUCache_GetMissingKey(t *testing.T) {
	c := NewLFUCache(3)
	if got := c.Get(42); got != -1 {
		t.Fatalf("Get(42) on empty cache = %d, want -1", got)
	}
}

func TestLFUCache_PutUpdatesValueAndBumpsFrequency(t *testing.T) {
	c := NewLFUCache(2)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Put(1, 100) // put on existing key = update + frequency bump
	c.Put(3, 3)   // key 2 (freq=1) is LFU, gets evicted, not key 1
	if got := c.Get(1); got != 100 {
		t.Fatalf("Get(1) = %d, want 100", got)
	}
	if got := c.Get(2); got != -1 {
		t.Fatalf("Get(2) = %d, want -1 (should have been evicted)", got)
	}
}

func TestLFUCache_TieBreakByLRUWithinFrequency(t *testing.T) {
	c := NewLFUCache(3)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Put(3, 3)
	// Touch 2 then 1, leaving 3 as least recently used among freq=1 ties... but here all
	// three still share freq=1 after their initial put, then we bump 1 and 2 via Get.
	c.Get(1) // freq(1)=2
	c.Get(2) // freq(2)=2
	// freq(3) is still 1, so it must be evicted first regardless of recency.
	c.Put(4, 4)
	if got := c.Get(3); got != -1 {
		t.Fatalf("Get(3) = %d, want -1 (lowest frequency must be evicted)", got)
	}
	if got := c.Get(1); got != 1 {
		t.Fatalf("Get(1) = %d, want 1", got)
	}
	if got := c.Get(2); got != 2 {
		t.Fatalf("Get(2) = %d, want 2", got)
	}
	if got := c.Get(4); got != 4 {
		t.Fatalf("Get(4) = %d, want 4", got)
	}
}

func TestLFUCache_MinFrequencyResetsAfterEviction(t *testing.T) {
	c := NewLFUCache(2)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Get(1)
	c.Get(1)
	c.Get(2) // freq: 1->3, 2->2
	c.Put(3, 3)
	// key 2 has the lowest freq now (2 vs 3), so it's evicted, not key 1.
	if got := c.Get(2); got != -1 {
		t.Fatalf("Get(2) = %d, want -1", got)
	}
	if got := c.Get(1); got != 1 {
		t.Fatalf("Get(1) = %d, want 1", got)
	}
	if got := c.Get(3); got != 3 {
		t.Fatalf("Get(3) = %d, want 3", got)
	}
}

func TestLFUCache_ReinsertAfterEviction(t *testing.T) {
	c := NewLFUCache(1)
	c.Put(1, 1)
	c.Put(2, 2) // evicts 1
	c.Put(1, 11)
	if got := c.Get(1); got != 11 {
		t.Fatalf("Get(1) = %d, want 11", got)
	}
	if got := c.Get(2); got != -1 {
		t.Fatalf("Get(2) = %d, want -1", got)
	}
}

func TestLFUCache_ZeroValueDistinguishableFromSentinel(t *testing.T) {
	c := NewLFUCache(2)
	c.Put(0, 0)
	if got := c.Get(0); got != 0 {
		t.Fatalf("Get(0) = %d, want 0 (must not be confused with the -1 miss sentinel)", got)
	}
}

func TestLFUCache_StressManyKeys(t *testing.T) {
	const capacity = 20
	c := NewLFUCache(capacity)

	for i := 0; i < capacity; i++ {
		c.Put(i, i)
	}
	// Keep keys 0..9 "hot" so they never get evicted.
	for round := 0; round < 5; round++ {
		for i := 0; i < 10; i++ {
			c.Get(i)
		}
	}
	// Insert new keys past capacity; only cold keys (10..19) should be evicted.
	for i := capacity; i < capacity+10; i++ {
		c.Put(i, i)
	}
	for i := 0; i < 10; i++ {
		if got := c.Get(i); got != i {
			t.Fatalf("hot key %d evicted unexpectedly: Get(%d) = %d, want %d", i, i, got, i)
		}
	}
}
