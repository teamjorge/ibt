package fifo

import (
	"testing"
)

func TestNewStore(t *testing.T) {
	t.Run("test NewStore with valid size", func(t *testing.T) {
		size := 5

		store := NewStore[int](size)

		if store.size != size {
			t.Errorf("expected store to have size %d. received: %d", size, store.size)
		}

		if store.l.size != size {
			t.Errorf("expected store's internal list to have size %d. received: %d", size, store.size)
		}

		if store.m == nil {
			t.Error("expected store's internal map to be initialized, but found nil")
		}
	})

	t.Run("test NewStore with 0 size", func(t *testing.T) {
		store := NewStore[int](-1)

		if store.size != storeDefaultSize {
			t.Errorf("expected store to have a default size %d. received: %d", storeDefaultSize, store.size)
		}

		if store.l.size != storeDefaultSize {
			t.Errorf("expected store's internal list to have a default size %d. received: %d", storeDefaultSize, store.size)
		}
	})
}

func TestStoreAdd(t *testing.T) {
	t.Run("test store add from empty", func(t *testing.T) {
		store := NewStore[int](3)
		store.Add(5)

		if store.current != 1 {
			t.Errorf("expected store to contain %d value(s). found %d", 1, store.current)
		}

		value, exists := store.m[5]
		if !exists {
			t.Errorf("expected item %d to be present in map and contain a non-nil node", 5)
		}

		if value.Key != 5 {
			t.Errorf("expected value .Key to be %d. received: %d", 5, value.Key)
		}
	})

	t.Run("test store add from existing", func(t *testing.T) {
		store := NewStore[int](3)
		store.Add(1)
		store.Add(2)

		if store.current != 2 {
			t.Errorf("expected store to contain %d value(s). found %d", 2, store.current)
		}

		value, exists := store.m[2]
		if !exists {
			t.Errorf("expected item %d to be present in map and contain a non-nil node", 2)
		}

		if value.Key != 2 {
			t.Errorf("expected value .Key to be %d. received: %d", 2, value.Key)
		}

		value, exists = store.m[1]
		if !exists {
			t.Errorf("expected item %d to be present in map and contain a non-nil node", 1)
		}

		if value.Key != 1 {
			t.Errorf("expected value .Key to be %d. received: %d", 1, value.Key)
		}
	})

	t.Run("test store add to capacity", func(t *testing.T) {
		store := NewStore[int](2)
		store.Add(1)
		store.Add(2)

		if store.current != 2 {
			t.Errorf("expected store to contain %d value(s). found %d", 2, store.current)
		}

		value, exists := store.m[2]
		if !exists {
			t.Errorf("expected item %d to be present in map and contain a non-nil node", 2)
		}

		if value.Key != 2 {
			t.Errorf("expected value .Key to be %d. received: %d", 2, value.Key)
		}

		value, exists = store.m[1]
		if !exists {
			t.Errorf("expected item %d to be present in map and contain a non-nil node", 1)
		}

		if value.Key != 1 {
			t.Errorf("expected value .Key to be %d. received: %d", 1, value.Key)
		}
	})

	t.Run("test store add to exceed capacity", func(t *testing.T) {
		store := NewStore[int](2)
		store.Add(1)
		store.Add(2)
		store.Add(3)
		store.Add(4)

		if store.current != 2 {
			t.Errorf("expected store to contain %d value(s). found %d", 2, store.current)
		}

		value, exists := store.m[3]
		if !exists {
			t.Errorf("expected item %d to be present in map and contain a non-nil node", 3)
		}

		if value.Key != 3 {
			t.Errorf("expected value .Key to be %d. received: %d", 3, value.Key)
		}

		value, exists = store.m[4]
		if !exists {
			t.Errorf("expected item %d to be present in map and contain a non-nil node", 4)
		}

		if value.Key != 4 {
			t.Errorf("expected value .Key to be %d. received: %d", 2, value.Key)
		}

		value, exists = store.m[2]
		if exists {
			t.Errorf("expected item %d to not be present in map and contain. found associated node %v", 2, value)
		}

		value, exists = store.m[1]
		if exists {
			t.Errorf("expected item %d to not be present in map and contain. found associated node %v", 1, value)
		}
	})

	t.Run("test store add to existing and move to back", func(t *testing.T) {
		store := NewStore[int](4)
		store.Add(1)
		store.Add(2)
		store.Add(3)
		store.Add(4)

		store.Add(2)

		if store.current != 4 {
			t.Errorf("expected store to contain %d value(s). found %d", 4, store.current)
		}

		if store.l.Last() != store.m[2] {
			t.Errorf("expected item %d to be last item in store to contain. found %v instead", 2, store.l.Last())
		}
	})
}

func TestStoreDelete(t *testing.T) {
	t.Run("test store delete from empty should not cause a panic", func(t *testing.T) {
		store := NewStore[int](5)

		store.Delete(3)
		if store.current != 0 {
			t.Errorf("expected store to be empty. found %d items", store.current)
		}
	})

	t.Run("test store delete the only item", func(t *testing.T) {
		store := NewStore[int](5)
		store.Add(1)

		store.Delete(1)

		if store.current != 0 {
			t.Errorf("expected store to be empty. found %d items", store.current)
		}

		value, exists := store.m[1]
		if exists {
			t.Errorf("expected item %d to not be present in map. found %+v", 1, value)
		}

		if store.l.Last() != nil {
			t.Errorf("expected last item in underlying store list to be empty. received %v", store.l.Last())
		}

		if store.l.First() != nil {
			t.Errorf("expected first item in underlying store list to be empty. received %v", store.l.First())
		}

		if store.l.Len() != 0 {
			t.Errorf("expected underlying store list to have a length of %d. received %d", 0, store.l.Len())
		}
	})

	t.Run("test store delete the middle item", func(t *testing.T) {
		store := NewStore[int](5)
		store.Add(1)
		store.Add(2)
		store.Add(3)

		store.Delete(2)

		if store.current != 2 {
			t.Errorf("expected store to be empty. found %d items", store.current)
		}

		value, exists := store.m[2]
		if exists {
			t.Errorf("expected item %d to not be present in map. found %+v", 2, value)
		}

		_, exists = store.m[1]
		if !exists {
			t.Errorf("expected item %d to be present in map", 1)
		}

		_, exists = store.m[3]
		if !exists {
			t.Errorf("expected item %d to be present in map", 3)
		}

		if store.l.Last() != store.m[3] {
			t.Errorf("expected last item in underlying store list to be key %d (%v). received %v",
				3, store.m[3], store.l.Last())
		}

		if store.l.First() != store.m[1] {
			t.Errorf("expected first item in underlying store list to be key %d (%v). received %v",
				1, store.m[1], store.l.First())
		}

		if store.l.Len() != 2 {
			t.Errorf("expected underlying store list to have a length of %d. received %d", 2, store.l.Len())
		}
	})
}

func TestStoreExists(t *testing.T) {
	t.Run("test store exists empty", func(t *testing.T) {
		store := NewStore[int](1)

		if store.Exists(5) {
			t.Errorf("expected %d to not exist in empty store. found value: %v", 5, store.m[5])
		}
	})

	t.Run("test store exists normal", func(t *testing.T) {
		store := NewStore[int](3)

		store.Add(3)
		store.Add(5)
		store.Add(7)

		if !store.Exists(5) {
			t.Errorf("expected %d to exist in store. currrent items %v", 5, store.m)
		}
	})

	t.Run("test store exists deleted item", func(t *testing.T) {
		store := NewStore[int](3)

		store.Add(3)
		store.Add(5)
		store.Add(7)
		store.Delete(5)

		if store.Exists(5) {
			t.Errorf("expected %d to not exist after being deleted. found value: %v", 5, store.m[5])
		}
	})
}

func TestStoreLen(t *testing.T) {
	t.Run("test store empty len", func(t *testing.T) {
		store := NewStore[int](5)
		if store.Len() != 0 {
			t.Errorf("expected empty store to have length of 0. found %d item. %v", store.Len(), store.m)
		}
	})

	t.Run("test store normal len", func(t *testing.T) {
		store := NewStore[int](5)
		store.Add(3)
		store.Add(5)
		store.Add(9)

		if store.Len() != 3 {
			t.Errorf("expected empty store to have length of %d. found %d item. %v", 3, store.Len(), store.m)
		}
	})

	t.Run("test store with deleted items len", func(t *testing.T) {
		store := NewStore[int](5)
		store.Add(3)
		store.Add(5)
		store.Add(9)
		store.Delete(5)
		store.Delete(3)

		if store.Len() != 1 {
			t.Errorf("expected empty store to have length of %d. found %d item. %v", 1, store.Len(), store.m)
		}
	})
}

/*************************Benchmarks*************************/

func BenchmarkStoreAdd(b *testing.B) {
	store := NewStore[int](5)

	// If N is not supplied, we default to 10000
	iter := b.N
	if iter == 1 {
		iter = 10000
	}

	for i := 0; i < iter; i++ {
		store.Add(i)
	}
}
