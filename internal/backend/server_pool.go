package backend

import "sync"

// ServerPool is a pool of backends.
type ServerPool struct {
	headBck *Backend
	currentBck *Backend
	mu sync.RWMutex
}

// NewServerPool creates a new server pool.
func NewServerPool() *ServerPool {
	return &ServerPool{
		headBck: &Backend{},
	}
}

// AddBackend adds a backend to the server pool.
func (s *ServerPool) AddBackend(b *Backend) {
	currBck := s.headBck
	for currBck.nextBck != nil {
		currBck = currBck.nextBck
	}
	currBck.nextBck = b
}

// GetNextBackend returns the next backend in the server pool.
func (s *ServerPool) GetNextBackend() *Backend {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.headBck == nil {
		return nil
	}

	var currBck *Backend
	if s.currentBck.nextBck == nil {
		currBck = s.headBck
	} else {
		currBck = s.currentBck.nextBck
	}
	
	for !currBck.GetState() {
		if currBck.nextBck == nil {
			currBck = s.headBck
		} else {
			currBck = currBck.nextBck
		}
	}

	return currBck
}