package health

import (
	"load-balancer/internal/backend"
	"net/http"
	"time"
)

type HealthChecker struct {
	pool *backend.ServerPool
	interval time.Duration
	timeout time.Duration
}

func NewHealthChecker(pool *backend.ServerPool, interval, timeout time.Duration) *HealthChecker {
	return &HealthChecker{
		pool: pool,
		interval: interval,
		timeout: timeout,
	}
}

func (c *HealthChecker) Start() {
	ticker := time.NewTicker(c.interval)

	go func ()  {
		for range ticker.C {
			c.CheckAllBackends()
		}
	}()
}

func (c *HealthChecker) CheckAllBackends() error {
	client := http.Client{
		Timeout: c.timeout,
	}

	currBck := c.pool.GetHeadBackend()
	for currBck != nil {
		resp, err := client.Get(currBck.URL.Host + "/health")
		if err != nil {
			currBck.SetState(false)
			continue
		}
		err = resp.Body.Close()
		if err != nil {
			return err
		}

		currBck.SetState(true)
		currBck = currBck.GetNextBackend()
	}

	return nil
}