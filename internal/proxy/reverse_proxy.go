package proxy

import (
	"load-balancer/internal/backend"
	"net/http"
	"net/http/httputil"
	"time"
)

type ReverseProxy struct {
	transport *http.Transport
}

// NewReverseProxy creates a new reverse proxy.
func NewReverseProxy() *ReverseProxy {
	return &ReverseProxy{
		transport: &http.Transport{
			MaxIdleConns: 100,
			IdleConnTimeout: 90 * time.Second,
			MaxIdleConnsPerHost: 10,
		},
	}
}

func (p *ReverseProxy) Serve(w *http.ResponseWriter, r *http.Request, bck *backend.Backend) {
	bck.AddConnection()
	defer bck.RemoveConnection()

	proxy := httputil.NewSingleHostReverseProxy(bck.URL)
	proxy.Transport = p.transport
	
	// Modify request
	proxy.Director = func(r *http.Request) {
		r.URL.Scheme = bck.URL.Scheme
		r.URL.Host = bck.URL.Host
		r.URL.Path = bck.URL.Path
		r.Host = bck.URL.Host

		if _, ok := r.Header["User-Agent"]; !ok {
			r.Header.Set("User-Agent", "")
		}
	}

	// Handle errors
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		bck.SetState(false)
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"error": "bad gateway}`))
	}
	
	start := time.Now()
	proxy.ServeHTTP(*w, r)
	bck.SetResponseTime(time.Since(start))
}