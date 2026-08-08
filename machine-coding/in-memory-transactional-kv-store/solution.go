package transactionalkvstore

/*
See README.md for the full problem statement, requirements, and rubric.

TODO:
 1. Add the fields you need on Store below. Hint: a "stack of layers," where each
    layer is a map[string]*string (nil = explicit delete marker), starting with one
    permanent-storage layer that's always present and never removed.
 2. Implement NewStore: start with just the single permanent-storage layer.
 3. Implement Get/Set/Delete/Begin/Commit/Rollback/InTransaction per the README and
    hints — the stack-of-layers approach makes nested transactions almost fall out for
    free once Get/Commit/Rollback correctly walk/merge/discard layers.
*/

type Store struct {
	// TODO: add your fields here
}

func NewStore() *Store {
	// TODO: implement
	return &Store{}
}

func (s *Store) Get(key string) (string, bool) {
	// TODO: implement
	return "", false
}

func (s *Store) Set(key, value string) {
	// TODO: implement
}

func (s *Store) Delete(key string) {
	// TODO: implement
}

func (s *Store) Begin() {
	// TODO: implement
}

func (s *Store) Commit() bool {
	// TODO: implement
	return false
}

func (s *Store) Rollback() bool {
	// TODO: implement
	return false
}

func (s *Store) InTransaction() bool {
	// TODO: implement
	return false
}
