package fifo

const (
	storeDefaultSize int = 10
)

// Store is a fixed-capacity FIFO ordered storage structure
type Store[k comparable] struct {
	m map[k]*Node[k]
	l *List[k]

	size    int
	current int
}

// NewStores creates a new instance of a Store.
//
//   - size: specifies the fixed amount ( > 0 ) of values the Store can hold.
//
// If size is equal to or less than 0, the default size (10) will be used
func NewStore[k comparable](size int) Store[k] {
	if size <= 0 {
		size = storeDefaultSize
	}

	store := Store[k]{
		m:    make(map[k]*Node[k]),
		l:    NewList[k](size),
		size: size,
	}

	return store
}

// Add an item to the store and determine if it is new.
//
// The item will be added at the last available position in the store. If the item is new
// and the store is at capacity, the first item will be removed to make way for the new item.
// If the store is at capacity and the item already exists, the item will simply be moved to
// the last available location.
//
// Returns true if the item is newly added or false if it is an existing item.
func (s *Store[k]) Add(item k) bool {
	foundItem, ok := s.m[item]
	if ok {
		s.l.Push(foundItem)

		return false
	}

	if s.current >= s.size {
		first := s.l.First()
		delete(s.m, first.Key)
		s.l.Delete(first)
		s.current--
	}

	node := &Node[k]{Key: item}
	s.l.Add(node)
	s.m[item] = node
	s.current++

	return true
}

// Delete an item from the Store
func (s *Store[k]) Delete(item k) {
	ptr, ok := s.m[item]
	if !ok {
		return
	}

	s.l.Delete(ptr)
	delete(s.m, item)
	s.current--
}

// Exists checks whether an item is present in the Store
func (s *Store[k]) Exists(item k) bool {
	_, ok := s.m[item]

	return ok
}

// Len provides the current length of the Store
func (s *Store[k]) Len() int { return s.current }
