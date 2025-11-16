package backend

import (
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type Backend struct {
	URL *url.URL
	state bool
	connections int64
	responseTime time.Duration
	nextBck *Backend
	mu sync.RWMutex
}

func NewBackend(stringUrl string) (*Backend, error) {
	u, err := url.Parse(stringUrl)
	if err != nil {
		return nil, err
	}
	return &Backend{
		URL: u,
		state: true,
	}, nil
}

func (b *Backend) SetState(state bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.state = state
}

func (b *Backend) GetState() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.state
}

func (b *Backend) AddConnection() {
	atomic.AddInt64(&b.connections, 1)
}

func (b *Backend) RemoveConnection() {
	atomic.AddInt64(&b.connections, -1)
}

func (b *Backend) GetResponseTime() time.Duration {
	return b.responseTime
}

func (b *Backend) SetResponseTime(responseTime time.Duration) {
	b.responseTime = responseTime
}

func (b *Backend) GetNextBackend() *Backend {
	return b.nextBck
}