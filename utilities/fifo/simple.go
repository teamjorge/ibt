package fifo

const (
	simpleDefaultSize int = 5
)

// Simple FIFO list implementation with size constraints.
type Simple[k comparable] struct {
	head int
	size int

	l []k
}

// NewSimple creates an instance of a Simple FIFO list
//
//   - size: specifies the fixed amount of values the Simple can hold
//
// If size is equal to or less than 0, the default size (5) will be used
func NewSimple[k comparable](size int) *Simple[k] {
	if size <= 0 {
		size = simpleDefaultSize
	}
	return &Simple[k]{size: size, l: make([]k, size)}
}

// Add a new item to the FIFO list
//
// Item is added into the last index. If the list is at capacity, the first
// item will be shifted out
func (s *Simple[k]) Add(item k) {
	s.l[s.head] = item
	s.head = (s.head + 1) % s.size
}

// Get the item at the provided offset.
//
// A nil value of type k will be returned in the following scenarios:
//
//   - Offset is less than zero
//   - Offset is greater than the size of the Simple
//   - A value does not exist in the specified index
func (s *Simple[k]) Get(offset int) k {
	var def k
	if offset < 0 || offset >= len(s.l) || offset >= s.size {
		return def
	}

	return s.l[(s.head+offset)%s.size]
}
