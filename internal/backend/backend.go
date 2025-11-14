package backend

import (
	"net/url"
	"sync"
	"time"
)

type Backend struct {
	URL *url.URL
	State bool
	Connections int64
	ResponseTime time.Duration
	mu sync.RWMutex
}

func NewBackend(stringUrl string) (*Backend, error) {
	u, err := url.Parse(stringUrl)
	if err != nil {
		return nil, err
	}
	return &Backend{
		URL: u,
		State: true,
	}, nil
}