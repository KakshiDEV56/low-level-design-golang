package pluggablecache

/*
See README.md for the full problem statement, requirements, and rubric.

EvictionPolicy is a fixed interface — do not change its shape.
*/

// EvictionPolicy decides which key to evict next. A Cache delegates all
// eviction-ordering decisions to one of these, so it never needs to know which
// specific algorithm (LRU, FIFO, LFU, ...) is in use.
type EvictionPolicy[K comparable] interface {
	// KeyAccessed is called whenever key is read or written in the cache.
	KeyAccessed(key K)
	// KeyRemoved is called when key is explicitly removed from the cache (not via
	// automatic eviction — Evict() itself is responsible for untracking the key it
	// returns).
	KeyRemoved(key K)
	// Evict returns the key that should be evicted next (and removes it from this
	// policy's own tracking), or ok=false if there's nothing to evict.
	Evict() (key K, ok bool)
}

// TODO: implement Cache[K, V].
//   - NewCache[K, V](capacity int, policy EvictionPolicy[K]) *Cache[K, V]
//   - Get(key K) (V, bool)
//   - Put(key K, value V)
//   - Remove(key K)
//   - Len() int
//
// Cache must not contain ANY policy-specific logic — it only calls the EvictionPolicy
// interface's methods. See README.md for the exact semantics of each method.
type Cache[K comparable, V any] struct {
	// TODO: add your fields here
}

func NewCache[K comparable, V any](capacity int, policy EvictionPolicy[K]) *Cache[K, V] {
	// TODO: implement
	return &Cache[K, V]{}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	// TODO: implement
	var zero V
	return zero, false
}

func (c *Cache[K, V]) Put(key K, value V) {
	// TODO: implement
}

func (c *Cache[K, V]) Remove(key K) {
	// TODO: implement
}

func (c *Cache[K, V]) Len() int {
	// TODO: implement
	return 0
}

// TODO: implement lruPolicy, satisfying EvictionPolicy[K]. See README.md hints
// (container/list.List + a map[K]*list.Element works well).
type lruPolicy[K comparable] struct {
	// TODO: add your fields here
}

func NewLRUPolicy[K comparable]() EvictionPolicy[K] {
	// TODO: implement
	return &lruPolicy[K]{}
}

func (p *lruPolicy[K]) KeyAccessed(key K) {
	// TODO: implement
}

func (p *lruPolicy[K]) KeyRemoved(key K) {
	// TODO: implement
}

func (p *lruPolicy[K]) Evict() (K, bool) {
	// TODO: implement
	var zero K
	return zero, false
}

// TODO: implement fifoPolicy, satisfying EvictionPolicy[K]. Unlike LRU, KeyAccessed on
// an ALREADY-TRACKED key must be a no-op — only a genuinely new key affects order.
type fifoPolicy[K comparable] struct {
	// TODO: add your fields here
}

func NewFIFOPolicy[K comparable]() EvictionPolicy[K] {
	// TODO: implement
	return &fifoPolicy[K]{}
}

func (p *fifoPolicy[K]) KeyAccessed(key K) {
	// TODO: implement
}

func (p *fifoPolicy[K]) KeyRemoved(key K) {
	// TODO: implement
}

func (p *fifoPolicy[K]) Evict() (K, bool) {
	// TODO: implement
	var zero K
	return zero, false
}
