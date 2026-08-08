package pluggablecache

import "testing"

// Both policies satisfy the same interface.
var (
	_ EvictionPolicy[string] = NewLRUPolicy[string]()
	_ EvictionPolicy[string] = NewFIFOPolicy[string]()
)

func TestCache_BasicGetPut(t *testing.T) {
	c := NewCache[string, int](3, NewLRUPolicy[string]())
	c.Put("a", 1)
	if got, ok := c.Get("a"); !ok || got != 1 {
		t.Fatalf("Get(a) = (%v, %v), want (1, true)", got, ok)
	}
}

func TestCache_GetMissingKey(t *testing.T) {
	c := NewCache[string, int](3, NewLRUPolicy[string]())
	if _, ok := c.Get("missing"); ok {
		t.Fatal("Get(missing) ok = true, want false")
	}
}

func TestCache_LenTracksSize(t *testing.T) {
	c := NewCache[string, int](5, NewLRUPolicy[string]())
	if c.Len() != 0 {
		t.Fatalf("Len() on empty cache = %d, want 0", c.Len())
	}
	c.Put("a", 1)
	c.Put("b", 2)
	if c.Len() != 2 {
		t.Fatalf("Len() = %d, want 2", c.Len())
	}
	c.Put("a", 100) // update, not a new entry
	if c.Len() != 2 {
		t.Fatalf("Len() after updating existing key = %d, want 2", c.Len())
	}
	c.Remove("a")
	if c.Len() != 1 {
		t.Fatalf("Len() after Remove = %d, want 1", c.Len())
	}
}

func TestCache_UpdateExistingKeyDoesNotEvict(t *testing.T) {
	c := NewCache[string, int](2, NewLRUPolicy[string]())
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("a", 100) // update; must not evict b
	if _, ok := c.Get("b"); !ok {
		t.Fatal("key b was evicted by an update to an existing key, want it still present")
	}
	if got, _ := c.Get("a"); got != 100 {
		t.Fatalf("Get(a) = %d, want 100", got)
	}
}

func TestCache_RemoveThenReinsert(t *testing.T) {
	c := NewCache[string, int](2, NewLRUPolicy[string]())
	c.Put("a", 1)
	c.Remove("a")
	if _, ok := c.Get("a"); ok {
		t.Fatal("Get(a) after Remove ok = true, want false")
	}
	c.Put("a", 2)
	if got, ok := c.Get("a"); !ok || got != 2 {
		t.Fatalf("Get(a) after reinsert = (%v, %v), want (2, true)", got, ok)
	}
}

// ---- LRU-specific behavior ----

func TestCache_LRU_EvictsLeastRecentlyUsed(t *testing.T) {
	c := NewCache[string, int](2, NewLRUPolicy[string]())
	c.Put("a", 1)
	c.Put("b", 2)
	c.Get("a")      // "a" is now more recently used than "b"
	c.Put("c", 3) // capacity exceeded; "b" should be evicted

	if _, ok := c.Get("b"); ok {
		t.Fatal("Get(b) ok = true, want false — b should have been LRU-evicted")
	}
	if _, ok := c.Get("a"); !ok {
		t.Fatal("Get(a) ok = false, want true — a was recently accessed, should survive")
	}
	if _, ok := c.Get("c"); !ok {
		t.Fatal("Get(c) ok = false, want true — c was just inserted")
	}
}

func TestCache_LRU_GetRefreshesOrder(t *testing.T) {
	c := NewCache[string, int](2, NewLRUPolicy[string]())
	c.Put("a", 1)
	c.Put("b", 2)
	c.Get("b") // touch b, making a the LRU
	c.Put("c", 3)

	if _, ok := c.Get("a"); ok {
		t.Fatal("Get(a) ok = true, want false — a was LRU and should have been evicted")
	}
	if _, ok := c.Get("b"); !ok {
		t.Fatal("Get(b) ok = false, want true — b was accessed after a")
	}
}

// ---- FIFO-specific behavior ----

func TestCache_FIFO_EvictsOldestInserted(t *testing.T) {
	c := NewCache[string, int](2, NewFIFOPolicy[string]())
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3) // "a" was inserted first, should be evicted

	if _, ok := c.Get("a"); ok {
		t.Fatal("Get(a) ok = true, want false — a was the first inserted, should be FIFO-evicted")
	}
	if _, ok := c.Get("b"); !ok {
		t.Fatal("Get(b) ok = false, want true")
	}
	if _, ok := c.Get("c"); !ok {
		t.Fatal("Get(c) ok = false, want true")
	}
}

func TestCache_FIFO_AccessDoesNotAffectEvictionOrder(t *testing.T) {
	c := NewCache[string, int](2, NewFIFOPolicy[string]())
	c.Put("a", 1)
	c.Put("b", 2)
	// Repeatedly access "a" — for LRU this would save it, but FIFO must ignore access
	// entirely and still evict based on insertion order.
	for i := 0; i < 5; i++ {
		c.Get("a")
	}
	c.Put("c", 3)

	if _, ok := c.Get("a"); ok {
		t.Fatal("Get(a) ok = true, want false — FIFO must evict by insertion order regardless of access count")
	}
}

// ---- Policy-level tests (bypassing Cache entirely) ----

func TestLRUPolicy_EvictOnEmptyReturnsFalse(t *testing.T) {
	p := NewLRUPolicy[string]()
	if _, ok := p.Evict(); ok {
		t.Fatal("Evict() on an empty LRU policy ok = true, want false")
	}
}

func TestFIFOPolicy_EvictOnEmptyReturnsFalse(t *testing.T) {
	p := NewFIFOPolicy[string]()
	if _, ok := p.Evict(); ok {
		t.Fatal("Evict() on an empty FIFO policy ok = true, want false")
	}
}

func TestLRUPolicy_KeyRemovedUntracksKey(t *testing.T) {
	p := NewLRUPolicy[string]()
	p.KeyAccessed("a")
	p.KeyAccessed("b")
	p.KeyRemoved("a")

	got, ok := p.Evict()
	if !ok || got != "b" {
		t.Fatalf("Evict() = (%v, %v), want (\"b\", true) — a was removed, only b should remain", got, ok)
	}
}

// ---- Swappability: identical operation sequence, different (correct) outcomes ----

func TestCache_SwappingPolicyChangesEvictionBehaviorOnly(t *testing.T) {
	run := func(policy EvictionPolicy[string]) (survivorAfterAccess string) {
		c := NewCache[string, int](2, policy)
		c.Put("a", 1)
		c.Put("b", 2)
		c.Get("a") // access "a"
		c.Put("c", 3)
		if _, ok := c.Get("a"); ok {
			return "a"
		}
		if _, ok := c.Get("b"); ok {
			return "b"
		}
		return ""
	}

	lruSurvivor := run(NewLRUPolicy[string]())
	fifoSurvivor := run(NewFIFOPolicy[string]())

	if lruSurvivor != "a" {
		t.Fatalf("with LRUPolicy, survivor = %q, want \"a\" (recently accessed)", lruSurvivor)
	}
	if fifoSurvivor != "b" {
		t.Fatalf("with FIFOPolicy, survivor = %q, want \"b\" (a was inserted first, access is irrelevant)", fifoSurvivor)
	}
}
