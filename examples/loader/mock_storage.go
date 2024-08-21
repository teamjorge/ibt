package main

// This is a mock external storage client.
//
// Think of it as a database, API, or external file
type storage struct {
	batchesLoaded int
}

// Simple constructor
func newStorage() *storage { return new(storage) }

func (s *storage) Connect() error { return nil }

func (s *storage) Exec(data []map[string]interface{}) error {
	s.batchesLoaded += len(data)
	return nil
}

func (s *storage) Close() error { return nil }

func (s *storage) Loaded() int { return s.batchesLoaded }
