package fifo

const (
	listDefaultSize int = 10
)

// List is a distinct value doubly-linked list implementation with size constraints
type List[k comparable] struct {
	size    int
	current int

	root Node[k]
}

// NewList creates a new instance of a List.
//
//   - size: specifies the fixed amount of values the List can hold
//
// If size is equal to or less than 0, the default size (10) will be used
func NewList[k comparable](size int) *List[k] {
	if size <= 0 {
		size = listDefaultSize
	}

	return &List[k]{size: size}
}

// Add a new item to the List
//
// Adding an existing Node will move it to the back of the list.
func (l *List[k]) Add(item *Node[k]) {
	// If item already exists
	if item.prev != nil && item.next != nil {
		l.Push(item)
		return
	} else {
		// If the list is currently empty
		if l.root.next == nil {
			l.root.next = item
			l.root.prev = item
			item.next = item
			item.prev = item
		} else {
			// Remove the first node if the list is full
			if l.current >= l.size {
				l.Delete(l.root.next)
			}

			l.root.prev.next = item
			item.prev = l.root.prev
			item.next = l.root.next
			l.root.next.prev = item
		}
	}

	// Add node to the back
	l.root.prev = item

	l.current++
}

// Delete an item from the List
func (l *List[k]) Delete(item *Node[k]) {
	if item.prev == nil && item.next == nil {
		// Invalid item
		return
	} else if l.root.prev == item && l.root.next == item {
		// Node is the only item
		l.root.prev = nil
		l.root.next = nil
	} else if l.root.prev == item {
		// Node is the last
		l.root.prev = l.root.prev.prev
		l.root.prev.next = l.root.next
		l.root.next.prev = l.root.prev
	} else if l.root.next == item {
		// Node is the first
		l.root.next = l.root.next.next
		l.root.next.prev = l.root.prev
		l.root.prev.next = l.root.next
	} else {
		// Node has a central position in list
		item.prev.next = item.next
		item.next.prev = item.prev
	}

	l.current--
	item.next = nil
	item.prev = nil
}

// First item of the List
//
// Returns nil if the List is empty
func (l *List[k]) First() *Node[k] {
	return l.root.next
}

// First item of the List
//
// Returns nil if the List is empty
func (l *List[k]) Last() *Node[k] {
	return l.root.prev
}

// Len returns the length of the List
func (l *List[k]) Len() int { return l.current }

// Pull item to the front of the List
func (l *List[k]) Pull(item *Node[k]) {
	if item.next == nil || item.prev == nil {
		// Invalid item
	} else if l.root.next == item {
		// Item is already first
	} else if l.root.prev == item {
		// Item is the last item
		l.root.next = item
		l.root.prev = item.prev
	} else {
		item.prev.next = item.next
		item.next.prev = item.prev

		l.root.prev.next = item
		l.root.next.prev = item
		item.prev = l.root.prev
		item.next = l.root.next

		l.root.next = item
	}
}

// Push item to the back of the List
func (l *List[k]) Push(item *Node[k]) {
	if item.next == nil || item.prev == nil {
		// Invalid item
	} else if l.root.prev == item {
		// Item is already last
		return
	} else if l.root.next == item {
		// Item is the first item
		l.root.prev = item
		l.root.next = item.next
	} else {
		item.prev.next = item.next
		item.next.prev = item.prev

		l.root.prev.next = item
		l.root.next.prev = item
		item.prev = l.root.prev
		item.next = l.root.next

		l.root.prev = item
	}
}

// Unpack the keys of all nodes in the list
func (l *List[k]) Unpack() []k {
	out := make([]k, 0)

	// List is empty
	if l.root.next == nil {
		return out
	}

	// Get the initial node
	initialNode := l.First()
	out = append(out, initialNode.Key)

	node := initialNode.next
	// Iterate till we hit the first node again
	for node != initialNode {
		out = append(out, node.Key)
		node = node.next
	}

	return out
}
