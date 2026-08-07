package designhashmap

import "testing"

func TestMyHashMap_Example1(t *testing.T) {
	m := NewMyHashMap()
	m.Put(1, 1)
	m.Put(2, 2)
	if got := m.Get(1); got != 1 {
		t.Fatalf("Get(1) = %d, want 1", got)
	}
	if got := m.Get(3); got != -1 {
		t.Fatalf("Get(3) = %d, want -1", got)
	}
	m.Put(2, 1)
	if got := m.Get(2); got != 1 {
		t.Fatalf("Get(2) = %d, want 1", got)
	}
	m.Remove(2)
	if got := m.Get(2); got != -1 {
		t.Fatalf("Get(2) after Remove = %d, want -1", got)
	}
}

func TestMyHashMap_GetOnEmptyMap(t *testing.T) {
	m := NewMyHashMap()
	if got := m.Get(0); got != -1 {
		t.Fatalf("Get(0) on empty map = %d, want -1", got)
	}
}

func TestMyHashMap_RemoveMissingKeyIsNoOp(t *testing.T) {
	m := NewMyHashMap()
	m.Put(1, 100)
	m.Remove(999) // key doesn't exist, must not panic or corrupt state
	if got := m.Get(1); got != 100 {
		t.Fatalf("Get(1) after removing unrelated key = %d, want 100", got)
	}
}

func TestMyHashMap_UpdateExistingKey(t *testing.T) {
	m := NewMyHashMap()
	m.Put(5, 50)
	m.Put(5, 500)
	m.Put(5, 5000)
	if got := m.Get(5); got != 5000 {
		t.Fatalf("Get(5) = %d, want 5000 (last Put should win)", got)
	}
}

func TestMyHashMap_ZeroValueIsDistinctFromMissing(t *testing.T) {
	m := NewMyHashMap()
	m.Put(7, 0)
	if got := m.Get(7); got != 0 {
		t.Fatalf("Get(7) = %d, want 0 (a stored 0 must not look like -1 sentinel)", got)
	}
}

func TestMyHashMap_CollisionsWithinSameBucket(t *testing.T) {
	// Keys that likely land in the same bucket under a simple `key % N` hash for common
	// bucket counts (1000, 769, ...): differ by a large multiple.
	m := NewMyHashMap()
	m.Put(1, 111)
	m.Put(1001, 222)
	m.Put(2001, 333)
	if got := m.Get(1); got != 111 {
		t.Fatalf("Get(1) = %d, want 111", got)
	}
	if got := m.Get(1001); got != 222 {
		t.Fatalf("Get(1001) = %d, want 222", got)
	}
	if got := m.Get(2001); got != 333 {
		t.Fatalf("Get(2001) = %d, want 333", got)
	}
	m.Remove(1001)
	if got := m.Get(1001); got != -1 {
		t.Fatalf("Get(1001) after Remove = %d, want -1", got)
	}
	// Removing the middle entry of a chain must not disturb its neighbors.
	if got := m.Get(1); got != 111 {
		t.Fatalf("Get(1) after removing 1001 = %d, want 111", got)
	}
	if got := m.Get(2001); got != 333 {
		t.Fatalf("Get(2001) after removing 1001 = %d, want 333", got)
	}
}

func TestMyHashMap_ReinsertAfterRemove(t *testing.T) {
	m := NewMyHashMap()
	m.Put(9, 90)
	m.Remove(9)
	m.Put(9, 900)
	if got := m.Get(9); got != 900 {
		t.Fatalf("Get(9) = %d, want 900", got)
	}
}

func TestMyHashMap_BoundaryKeysAndValues(t *testing.T) {
	m := NewMyHashMap()
	m.Put(0, 0)
	m.Put(1000000, 1000000)
	if got := m.Get(0); got != 0 {
		t.Fatalf("Get(0) = %d, want 0", got)
	}
	if got := m.Get(1000000); got != 1000000 {
		t.Fatalf("Get(1000000) = %d, want 1000000", got)
	}
}

func TestMyHashMap_StressManyKeys(t *testing.T) {
	m := NewMyHashMap()
	const n = 5000

	for i := 0; i < n; i++ {
		m.Put(i, i*2)
	}
	for i := 0; i < n; i++ {
		if got := m.Get(i); got != i*2 {
			t.Fatalf("Get(%d) = %d, want %d", i, got, i*2)
		}
	}
	for i := 0; i < n; i += 2 {
		m.Remove(i)
	}
	for i := 0; i < n; i++ {
		want := i * 2
		if i%2 == 0 {
			want = -1
		}
		if got := m.Get(i); got != want {
			t.Fatalf("Get(%d) after removing evens = %d, want %d", i, got, want)
		}
	}
}
