package fifo

import (
	"testing"
)

func TestNewList(t *testing.T) {
	t.Run("test NewList with valid size", func(t *testing.T) {
		size := 5

		list := NewList[int](size)

		if list.size != size {
			t.Errorf("expected list to have size %d. received: %d", size, list.size)
		}

		if list.current != 0 {
			t.Errorf("expected new list to have a current value of 0. received: %d", list.current)
		}
	})

	t.Run("test NewList with 0 size", func(t *testing.T) {
		list := NewList[int](-1)

		if list.size != listDefaultSize {
			t.Errorf("expected list to have a default size %d. received: %d", listDefaultSize, list.size)
		}

		if list.current != 0 {
			t.Errorf("expected new list to have a current value of 0. received: %d", list.current)
		}
	})
}

func TestListAdd(t *testing.T) {
	t.Run("test list add from empty", func(t *testing.T) {
		list := NewList[int](3)

		node := &Node[int]{Key: 5}

		list.Add(node)

		if list.current != 1 {
			t.Errorf("expected list to contain %d value(s). found %d", 1, list.current)
		}

		if list.root.prev != node && list.root.next != node {
			t.Errorf("expected node to be both the first and last node. root node: %v", list.root)
		}
	})

	t.Run("test list add from existing", func(t *testing.T) {
		list := NewList[int](3)
		node1 := &Node[int]{Key: 1}
		node2 := &Node[int]{Key: 2}

		list.Add(node1)
		list.Add(node2)

		if list.current != 2 {
			t.Errorf("expected list to contain %d value(s). found %d", 2, list.current)
		}

		if list.root.next != node1 {
			t.Errorf("expected node1 to be first node. first node: %v", list.root.next)
		}

		if list.root.prev != node2 {
			t.Errorf("expected node2 to be last node. last node: %v", list.root.prev)
		}

		if node1.next != node2 {
			t.Errorf("expected the node after node1 to be node2. found: %v", node1.next)
		}
		if node1.prev != node2 {
			t.Errorf("expected the node before node1 to be node2. found: %v", node1.prev)
		}

		if node2.prev != node1 {
			t.Errorf("expected the node after node2 to be node1. found: %v", node2.next)
		}
		if node2.next != node1 {
			t.Errorf("expected the node before node2 to be node1. found: %v", node2.prev)
		}
	})

	t.Run("test list add to capacity", func(t *testing.T) {
		list := NewList[int](2)
		node1 := &Node[int]{Key: 1}
		node2 := &Node[int]{Key: 2}

		list.Add(node1)
		list.Add(node2)

		if list.current != 2 {
			t.Errorf("expected list to contain %d value(s). found %d", 2, list.current)
		}

		if list.root.next != node1 {
			t.Errorf("expected node1 to be first node. first node: %v", list.root.next)
		}

		if list.root.prev != node2 {
			t.Errorf("expected node2 to be last node. last node: %v", list.root.prev)
		}

		if node1.next != node2 {
			t.Errorf("expected the node after node1 to be node2. found: %v", node1.next)
		}
		if node1.prev != node2 {
			t.Errorf("expected the node before node1 to be node2. found: %v", node1.prev)
		}

		if node2.prev != node1 {
			t.Errorf("expected the node after node2 to be node1. found: %v", node2.next)
		}
		if node2.next != node1 {
			t.Errorf("expected the node before node2 to be node1. found: %v", node2.prev)
		}
	})

	t.Run("test list add to exceed capacity", func(t *testing.T) {
		list := NewList[int](2)
		node1 := &Node[int]{Key: 1}
		node2 := &Node[int]{Key: 2}
		node3 := &Node[int]{Key: 3}

		list.Add(node1)
		list.Add(node2)
		list.Add(node3)

		if list.current != 2 {
			t.Errorf("expected list to contain %d value(s). found %d", 2, list.current)
		}

		if list.root.next != node2 {
			t.Errorf("expected node2 to be first node. first node: %v", list.root.next)
		}

		if list.root.prev != node3 {
			t.Errorf("expected node3 to be last node. last node: %v", list.root.prev)
		}

		if node2.next != node3 {
			t.Errorf("expected the node after node2 to be node3. found: %v", node1.next)
		}
		if node2.prev != node3 {
			t.Errorf("expected the node before node2 to be node3. found: %v", node1.prev)
		}

		if node3.prev != node2 {
			t.Errorf("expected the node after node3 to be node2. found: %v", node2.next)
		}
		if node3.next != node2 {
			t.Errorf("expected the node before node3 to be node2. found: %v", node2.prev)
		}
	})
}

func TestListDelete(t *testing.T) {
	t.Run("test list delete empty", func(t *testing.T) {
		list := NewList[int](3)

		list.Delete(&Node[int]{Key: 2})

		if list.current != 0 {
			t.Errorf("expected list to be empty. found %d items", list.current)
		}
	})

	t.Run("test list delete single item", func(t *testing.T) {
		list := NewList[int](3)
		node1 := &Node[int]{Key: 1}

		list.Add(node1)
		list.Delete(node1)

		if list.current != 0 {
			t.Errorf("expected list to be empty. found %d items", list.current)
		}

		if list.root.prev != nil {
			t.Errorf("expected root prev to be nil. found %v instead", list.root.prev)
		}

		if list.root.next != nil {
			t.Errorf("expected root next to be nil. found %v instead", list.root.next)
		}
	})

	t.Run("test list delete full", func(t *testing.T) {
		list := NewList[int](3)
		node1 := &Node[int]{Key: 1}
		node2 := &Node[int]{Key: 2}
		node3 := &Node[int]{Key: 3}

		list.Add(node1)
		list.Add(node2)
		list.Add(node3)

		list.Delete(node3)

		if list.current != 2 {
			t.Errorf("expected list to have %d items. found %d items", 2, list.current)
		}

		if list.root.prev != node2 {
			t.Errorf("expected node2 to be last node. first node: %v", list.root.prev)
		}

		if list.root.next != node1 {
			t.Errorf("expected node1 to be first node. last node: %v", list.root.next)
		}

		if node1.next != node2 {
			t.Errorf("expected the node after node1 to be node2. found: %v", node1.next)
		}
		if node1.prev != node2 {
			t.Errorf("expected the node before node1 to be node2. found: %v", node1.prev)
		}

		if node2.prev != node1 {
			t.Errorf("expected the node after node2 to be node1. found: %v", node2.next)
		}
		if node2.next != node1 {
			t.Errorf("expected the node before node2 to be node1. found: %v", node2.prev)
		}
	})
}

func TestListFirst(t *testing.T) {
	t.Run("test list first empty", func(t *testing.T) {
		list := NewList[int](3)

		if list.First() != nil {
			t.Errorf("expected empty list to not have a first item. found %v", list.First())
		}
	})

	t.Run("test list first non-full list", func(t *testing.T) {
		list := NewList[int](4)
		node1 := &Node[int]{Key: 1}
		node2 := &Node[int]{Key: 2}
		node3 := &Node[int]{Key: 3}

		list.Add(node1)
		list.Add(node2)
		list.Add(node3)

		if list.First() != node1 {
			t.Errorf("expected node1 to be the first item. found %v", list.First())
		}

		if node1.next != node2 {
			t.Errorf("expected node2 to be the next item. found %v", node1.next)
		}

		if node1.prev != node3 {
			t.Errorf("expected node3 to be the previous item. found %v", node1.prev)
		}
	})

	t.Run("test list first full list", func(t *testing.T) {
		list := NewList[int](3)
		node1 := &Node[int]{Key: 1}
		node2 := &Node[int]{Key: 2}
		node3 := &Node[int]{Key: 3}
		node4 := &Node[int]{Key: 4}

		list.Add(node1)
		list.Add(node2)
		list.Add(node3)
		list.Add(node4)

		if list.First() != node2 {
			t.Errorf("expected node2 to be the first item. found %v", list.First())
		}

		if node2.next != node3 {
			t.Errorf("expected node3 to be the next item. found %v", node1.next)
		}

		if node2.prev != node4 {
			t.Errorf("expected node3 to be the previous item. found %v", node1.prev)
		}
	})
}

func TestListLast(t *testing.T) {
	t.Run("test list last empty", func(t *testing.T) {
		list := NewList[int](3)

		if list.Last() != nil {
			t.Errorf("expected empty list to not have a last item. found %v", list.Last())
		}
	})

	t.Run("test list last non-full list", func(t *testing.T) {
		list := NewList[int](4)
		node1 := &Node[int]{Key: 1}
		node2 := &Node[int]{Key: 2}
		node3 := &Node[int]{Key: 3}

		list.Add(node1)
		list.Add(node2)
		list.Add(node3)

		if list.Last() != node3 {
			t.Errorf("expected node3 to be the last item. found %v", list.Last())
		}

		if node3.next != node1 {
			t.Errorf("expected node1 to be the next item. found %v", node3.next)
		}

		if node3.prev != node2 {
			t.Errorf("expected node2 to be the previous item. found %v", node3.prev)
		}
	})

	t.Run("test list first full list", func(t *testing.T) {
		list := NewList[int](3)
		node1 := &Node[int]{Key: 1}
		node2 := &Node[int]{Key: 2}
		node3 := &Node[int]{Key: 3}
		node4 := &Node[int]{Key: 4}

		list.Add(node1)
		list.Add(node2)
		list.Add(node3)
		list.Add(node4)

		if list.Last() != node4 {
			t.Errorf("expected node4 to be the last item. found %v", list.Last())
		}

		if node4.next != node2 {
			t.Errorf("expected node2 to be the next item. found %v", node4.next)
		}

		if node4.prev != node3 {
			t.Errorf("expected node3 to be the previous item. found %v", node4.prev)
		}
	})
}

func TestListPull(t *testing.T) {
	t.Run("test pull empty list does not panic", func(t *testing.T) {
		list := NewList[int](4)

		list.Pull(&Node[int]{Key: 1})

		if list.current != 0 {
			t.Errorf("expected list current to be 0. received: %d", list.current)
		}
	})

	t.Run("test pull single item list", func(t *testing.T) {
		list := NewList[int](4)
		node1 := &Node[int]{Key: 1}

		list.Add(node1)
		list.Pull(node1)

		if list.current != 1 {
			t.Errorf("expected list current to be 1. received: %d", list.current)
		}

		if list.root.prev != node1 && list.root.next != node1 {
			t.Errorf("expected node to be both the first and last node. root node: %v", list.root)
		}
	})

	t.Run("test pull non-full list", func(t *testing.T) {
		list := NewList[int](4)
		node1 := &Node[int]{Key: 1}
		node2 := &Node[int]{Key: 2}
		node3 := &Node[int]{Key: 3}

		list.Add(node1)
		list.Add(node2)
		list.Add(node3)
		list.Pull(node3)

		if list.current != 3 {
			t.Errorf("expected list current to be 1. received: %d", list.current)
		}

		if list.First() != node3 {
			t.Errorf("expected node3 to be the first item. found %v", list.First())
		}

		if list.Last() != node2 {
			t.Errorf("expected node2 to be the last item. found %v", list.Last())
		}

		if node3.next != node1 {
			t.Errorf("expected node1 to be the next item. found %v", node3.next)
		}

		if node3.prev != node2 {
			t.Errorf("expected node2 to be the previous item. found %v", node3.prev)
		}
	})

	t.Run("test pull full list", func(t *testing.T) {
		list := NewList[int](4)
		node1 := &Node[int]{Key: 1}
		node2 := &Node[int]{Key: 2}
		node3 := &Node[int]{Key: 3}
		node4 := &Node[int]{Key: 4}

		list.Add(node1)
		list.Add(node2)
		list.Add(node3)
		list.Add(node4)
		list.Pull(node3)

		if list.current != 4 {
			t.Errorf("expected list current to be 1. received: %d", list.current)
		}

		if list.First() != node3 {
			t.Errorf("expected node3 to be the first item. found %v", list.First())
		}

		if list.Last() != node4 {
			t.Errorf("expected node4 to be the last item. found %v", list.Last())
		}

		if node3.next != node1 {
			t.Errorf("expected node3 to be the next item. found %v", node3.next)
		}

		if node3.prev != node4 {
			t.Errorf("expected node4 to be the previous item. found %v", node3.prev)
		}
	})
}

func TestListPush(t *testing.T) {
	t.Run("test push empty list does not panic", func(t *testing.T) {
		list := NewList[int](4)

		list.Push(&Node[int]{Key: 1})

		if list.current != 0 {
			t.Errorf("expected list current to be 0. received: %d", list.current)
		}
	})

	t.Run("test push single item list", func(t *testing.T) {
		list := NewList[int](4)
		node1 := &Node[int]{Key: 1}

		list.Add(node1)
		list.Push(node1)

		if list.current != 1 {
			t.Errorf("expected list current to be 1. received: %d", list.current)
		}

		if list.root.prev != node1 && list.root.next != node1 {
			t.Errorf("expected node to be both the first and last node. root node: %v", list.root)
		}
	})

	t.Run("test push non-full list", func(t *testing.T) {
		list := NewList[int](4)
		node1 := &Node[int]{Key: 1}
		node2 := &Node[int]{Key: 2}
		node3 := &Node[int]{Key: 3}

		list.Add(node1)
		list.Add(node2)
		list.Add(node3)
		list.Push(node1)

		if list.current != 3 {
			t.Errorf("expected list current to be 1. received: %d", list.current)
		}

		if list.Last() != node1 {
			t.Errorf("expected node1 to be the last item. found %v", list.Last())
		}

		if list.First() != node2 {
			t.Errorf("expected node2 to be the first item. found %v", list.First())
		}

		if node1.next != node2 {
			t.Errorf("expected node2 to be the next item. found %v", node1.next)
		}

		if node1.prev != node3 {
			t.Errorf("expected node3 to be the previous item. found %v", node1.prev)
		}
	})

	t.Run("test push full list", func(t *testing.T) {
		list := NewList[int](4)
		node1 := &Node[int]{Key: 1}
		node2 := &Node[int]{Key: 2}
		node3 := &Node[int]{Key: 3}
		node4 := &Node[int]{Key: 4}

		list.Add(node1)
		list.Add(node2)
		list.Add(node3)
		list.Add(node4)
		list.Push(node3)

		if list.current != 4 {
			t.Errorf("expected list current to be 1. received: %d", list.current)
		}

		if list.First() != node1 {
			t.Errorf("expected node1 to be the last item. found %v", list.First())
		}

		if list.Last() != node3 {
			t.Errorf("expected node3 to be the last item. found %v", list.Last())
		}

		if node3.next != node1 {
			t.Errorf("expected node1 to be the next item. found %v", node4.next)
		}

		if node3.prev != node4 {
			t.Errorf("expected node4 to be the previous item. found %v", node4.prev)
		}
	})
}

func TestListUnpack(t *testing.T) {
	t.Run("test unpack empty list", func(t *testing.T) {
		list := NewList[int](5)

		unpacked := list.Unpack()

		if len(unpacked) != 0 {
			t.Errorf("expected unpacked len to be 0. received: %d", len(unpacked))
		}
	})

	t.Run("test unpack single item list", func(t *testing.T) {
		list := NewList[int](5)

		list.Add(&Node[int]{Key: 3})

		unpacked := list.Unpack()

		if len(unpacked) != 1 {
			t.Errorf("expected unpacked len to be 1. received: %d", len(unpacked))
		}

		if unpacked[0] != 3 {
			t.Errorf("expected unpacked item to be %d. received: %d", unpacked[0], len(unpacked))
		}
	})

	t.Run("test unpack full item list", func(t *testing.T) {
		list := NewList[int](5)

		list.Add(&Node[int]{Key: 3})
		list.Add(&Node[int]{Key: 4})
		list.Add(&Node[int]{Key: 5})
		list.Add(&Node[int]{Key: 6})
		list.Add(&Node[int]{Key: 7})

		unpacked := list.Unpack()

		if len(unpacked) != 5 {
			t.Errorf("expected unpacked len to be 5. received: %d", len(unpacked))
		}

		if unpacked[0] != 3 {
			t.Errorf("expected unpacked item to be %d. received: %d", unpacked[0], len(unpacked))
		}

		if unpacked[2] != 5 {
			t.Errorf("expected unpacked item to be %d. received: %d", unpacked[2], len(unpacked))
		}

		if unpacked[4] != 7 {
			t.Errorf("expected unpacked item to be %d. received: %d", unpacked[4], len(unpacked))
		}
	})

	t.Run("test unpack correct order", func(t *testing.T) {
		list := NewList[int](3)

		node1 := &Node[int]{Key: 1}
		node2 := &Node[int]{Key: 2}
		node3 := &Node[int]{Key: 3}
		node4 := &Node[int]{Key: 4}
		node5 := &Node[int]{Key: 5}
		node6 := &Node[int]{Key: 6}
		node7 := &Node[int]{Key: 7}

		list.Add(node1)
		list.Add(node2)
		list.Add(node3)

		//1-2-3

		list.Push(node1)
		list.Delete(node2)

		//3-1

		list.Add(node7)
		list.Delete(node7)

		//3-1-7
		//3-1

		list.Add(node1)
		list.Add(node5)

		//1-1-5
		list.Add(node4)
		list.Add(node6)

		// 5-4-6

		unpacked := list.Unpack()

		if len(unpacked) != 3 {
			t.Errorf("expected unpacked len to be 3. received: %d", len(unpacked))
		}

		if unpacked[0] != 5 {
			t.Errorf("expected unpacked item to be %d. received: %d", 5, unpacked[0])
		}

		if unpacked[1] != 4 {
			t.Errorf("expected unpacked item to be %d. received: %d", 4, unpacked[0])
		}

		if unpacked[2] != 6 {
			t.Errorf("expected unpacked item to be %d. received: %d", 6, unpacked[2])
		}
	})
}
