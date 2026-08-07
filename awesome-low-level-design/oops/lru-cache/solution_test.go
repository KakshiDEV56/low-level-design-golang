package lrucache

import "testing"

func TestLRUCache_Example1(t *testing.T) {
	c := NewLRUCache(2)
	c.Put(1, 1)
	c.Put(2, 2)
	if got := c.Get(1); got != 1 {
		t.Fatalf("Get(1) = %d, want 1", got)
	}
	c.Put(3, 3) // evicts key 2 (least recently used)
	if got := c.Get(2); got != -1 {
		t.Fatalf("Get(2) = %d, want -1", got)
	}
	c.Put(4, 4) // evicts key 1
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

func TestLRUCache_CapacityOne(t *testing.T) {
	c := NewLRUCache(1)
	c.Put(1, 10)
	if got := c.Get(1); got != 10 {
		t.Fatalf("Get(1) = %d, want 10", got)
	}
	c.Put(2, 20) // evicts key 1 immediately since capacity is 1
	if got := c.Get(1); got != -1 {
		t.Fatalf("Get(1) = %d, want -1", got)
	}
	if got := c.Get(2); got != 20 {
		t.Fatalf("Get(2) = %d, want 20", got)
	}
}

func TestLRUCache_GetMissingKey(t *testing.T) {
	c := NewLRUCache(2)
	if got := c.Get(99); got != -1 {
		t.Fatalf("Get(99) on empty cache = %d, want -1", got)
	}
}

func TestLRUCache_UpdateExistingKeyRefreshesRecency(t *testing.T) {
	c := NewLRUCache(2)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Put(1, 100) // updating key 1 should make it most recently used
	c.Put(3, 3)   // must evict key 2, not key 1
	if got := c.Get(1); got != 100 {
		t.Fatalf("Get(1) = %d, want 100", got)
	}
	if got := c.Get(2); got != -1 {
		t.Fatalf("Get(2) = %d, want -1 (should have been evicted)", got)
	}
	if got := c.Get(3); got != 3 {
		t.Fatalf("Get(3) = %d, want 3", got)
	}
}

func TestLRUCache_GetRefreshesRecency(t *testing.T) {
	c := NewLRUCache(2)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Get(1)    // touching key 1 makes it most recently used
	c.Put(3, 3) // must evict key 2, not key 1
	if got := c.Get(2); got != -1 {
		t.Fatalf("Get(2) = %d, want -1 (should have been evicted)", got)
	}
	if got := c.Get(1); got != 1 {
		t.Fatalf("Get(1) = %d, want 1", got)
	}
}

func TestLRUCache_PutExistingKeyDoesNotGrowSize(t *testing.T) {
	c := NewLRUCache(2)
	c.Put(1, 1)
	c.Put(1, 2)
	c.Put(1, 3)
	c.Put(2, 20)
	// Cache should still only hold {1,2}; no eviction should have happened.
	if got := c.Get(1); got != 3 {
		t.Fatalf("Get(1) = %d, want 3", got)
	}
	if got := c.Get(2); got != 20 {
		t.Fatalf("Get(2) = %d, want 20", got)
	}
}

func TestLRUCache_LargerCapacityManyOperations(t *testing.T) {
	c := NewLRUCache(3)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Put(3, 3)
	c.Get(1)    // order of recency: 2 (LRU), 3, 1 (MRU)
	c.Put(4, 4) // evicts 2
	if got := c.Get(2); got != -1 {
		t.Fatalf("Get(2) = %d, want -1", got)
	}
	if got := c.Get(3); got != 3 {
		t.Fatalf("Get(3) = %d, want 3", got)
	}
	if got := c.Get(1); got != 1 {
		t.Fatalf("Get(1) = %d, want 1", got)
	}
	if got := c.Get(4); got != 4 {
		t.Fatalf("Get(4) = %d, want 4", got)
	}
}

func TestLRUCache_ZeroAndBoundaryValues(t *testing.T) {
	c := NewLRUCache(2)
	c.Put(0, 0)
	if got := c.Get(0); got != 0 {
		t.Fatalf("Get(0) = %d, want 0 (zero value must be distinguishable from -1 sentinel)", got)
	}
}

func TestLRUCache_StressSequential(t *testing.T) {
	const capacity = 50
	c := NewLRUCache(capacity)

	for i := 0; i < 500; i++ {
		c.Put(i, i*10)
	}

	// Only the last `capacity` keys (450..499) should remain.
	for i := 0; i < 500-capacity; i++ {
		if got := c.Get(i); got != -1 {
			t.Fatalf("Get(%d) = %d, want -1 (should have been evicted)", i, got)
		}
	}
	for i := 500 - capacity; i < 500; i++ {
		if got := c.Get(i); got != i*10 {
			t.Fatalf("Get(%d) = %d, want %d", i, got, i*10)
		}
	}
}

func TestLRUCache_RepeatedGetsDoNotEvict(t *testing.T) {
	c := NewLRUCache(2)
	c.Put(1, 1)
	c.Put(2, 2)
	for i := 0; i < 10; i++ {
		c.Get(1)
		c.Get(2)
	}
	if got := c.Get(1); got != 1 {
		t.Fatalf("Get(1) = %d, want 1", got)
	}
	if got := c.Get(2); got != 2 {
		t.Fatalf("Get(2) = %d, want 2", got)
	}
}
