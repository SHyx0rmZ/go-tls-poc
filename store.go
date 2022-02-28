package tls

import (
	"sync"
)

// CleanupFunc is used for functions that will remove a
// go-routine's TLS, so as to not leak memory when leaving
// a go-routine's scope.
type CleanupFunc func()

// Close will remove a go-routine's TLS. It is only provided
// as a convenience, you can also just call the CleanupFunc
// directly.
func (fn CleanupFunc) Close() error {
	fn()
	return nil
}

// Store is a type that can be used to create a TLS that can
// store values of a certain data type. This should usually be
// something like a struct or a map which actually store the
// data you are interested in. Store should not be used to
// create many TLS with different data types.
type Store struct {
	mu sync.Mutex
	m  map[int64]interface{}

	// New shall return a new object which will be the base
	// of each go-routine's freshly created TLS. This usually
	// should be a struct or map which contains the actual data
	// you want to store per go-routine.
	New func() interface{}
}

// Delete removes the current go-routine's TLS. This can only
// be called from the go-routine that created the TLS. To avoid
// leaks in situations where you can't call Delete from the same
// go-routine call Closer after creating the TLS and call the
// returned CleanupFunc when you're done with the TLS.
func (s *Store) Delete() {
	s.delete(goid())
}

// Load returns the current go-routine's TLS. If it doesn't
// exist yet, New will be called to create a new base value.
func (s *Store) Load() interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	goid := goid()
	v, ok := s.m[goid]
	if !ok {
		v = s.New()
		if s.m == nil {
			s.m = make(map[int64]interface{})
		}
		s.m[goid] = v
	}
	return v
}

// Closer returns a CleanupFunc that can be used to remove a
// go-routine's TLS. If you can you should prefer using Delete
// over using this function.
func (s *Store) Closer() CleanupFunc {
	goid := goid()
	return func() {
		s.delete(goid)
	}
}

// Store updates the current go-routine's TLS. This is only
// necessary when the supplied New doesn't return a pointer.
func (s *Store) Store(v interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.m == nil {
		s.m = make(map[int64]interface{})
	}
	s.m[goid()] = v
}

func (s *Store) delete(goid int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, goid)
}
