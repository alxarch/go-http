package vhosts

import (
	"net/http"
	"sync"
)

type VirtualHosts struct {
	vhosts         map[string]http.Handler
	mu             sync.RWMutex
	defaultHandler http.Handler
}

func NewVirtualHosts(vhosts map[string]http.Handler) *VirtualHosts {
	v := &VirtualHosts{}
	for host, h := range vhosts {
		v.HandleHost(h, host)
	}
	return v
}

func (v *VirtualHosts) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	if v.vhosts != nil {
		if handler := v.vhosts[r.Host]; handler != nil {
			handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

func (v *VirtualHosts) HandleHost(handler http.Handler, hosts ...string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.vhosts == nil {
		v.vhosts = make(map[string]http.Handler)
	}
	for _, host := range hosts {
		v.vhosts[host] = handler
	}
}
