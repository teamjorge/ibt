package fifo

import "testing"

func TestNewSimple(t *testing.T) {
	t.Run("test invalid Simple size", func(t *testing.T) {
		simple := NewSimple[int](0)

		if len(simple.l) != simpleDefaultSize {
			t.Errorf("expected internal simple list length of %d. found %d", simpleDefaultSize, len(simple.l))
		}

		if simple.size != simpleDefaultSize {
			t.Errorf("expected internal simple list length of %d. found %d", simpleDefaultSize, len(simple.l))
		}
	})

	t.Run("test valid Simple size", func(t *testing.T) {
		simple := NewSimple[int](3)

		if len(simple.l) != 3 {
			t.Errorf("expected internal simple list length of %d. found %d", 3, len(simple.l))
		}

		if simple.size != 3 {
			t.Errorf("expected internal simple list length of %d. found %d", 3, len(simple.l))
		}
	})
}

func TestSimpleAdd(t *testing.T) {
	t.Run("test add to Simple from empty", func(t *testing.T) {
		simple := NewSimple[int](3)

		simple.Add(9)
		simple.Add(55)

		if simple.head != 2 {
			t.Errorf("expected simple head to be %d. found: %d", 2, simple.head)
		}

		if simple.l[0] != 9 {
			t.Errorf("expected first index to be %d. found: %d", 9, simple.l[0])
		}
	})

	t.Run("test add to Simple to over capacity", func(t *testing.T) {
		simple := NewSimple[int](3)

		simple.Add(9)
		simple.Add(55)
		simple.Add(69)
		simple.Add(93)

		if simple.head != 1 {
			t.Errorf("expected simple head to be %d. found: %d", 2, simple.head)
		}

		if simple.l[0] != 93 {
			t.Errorf("expected first index to be %d. found: %d", 93, simple.l[0])
		}

		if simple.l[1] != 55 {
			t.Errorf("expected second index to be %d. found: %d", 55, simple.l[1])
		}

		if simple.l[2] != 69 {
			t.Errorf("expected last index to be %d. found: %d", 69, simple.l[2])
		}
	})
}

func TestSimpleGet(t *testing.T) {
	t.Run("test get on empty simple", func(t *testing.T) {
		simple := NewSimple[int](3)

		if simple.Get(2) != 0 {
			t.Errorf("expected get on empty Simple to return nil value. received %d", simple.Get(4))
		}
	})

	t.Run("test get invalid offsets", func(t *testing.T) {
		simple := NewSimple[int](3)

		if simple.Get(-1) != 0 {
			t.Errorf("expected get with offset less than zero to return nil value. received %d", simple.Get(-1))
		}

		if simple.Get(5) != 0 {
			t.Errorf("expected get with offset greater than size to return nil value. received %d", simple.Get(5))
		}
	})
}
