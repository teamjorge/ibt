package fifo

// Node is used to track an element in a Store and it's underlying doubly-linked list.
type Node[K comparable] struct {
	next, prev *Node[K]

	// Key represents the actual value of a stored item.
	Key K
}
