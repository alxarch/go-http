package vhosts

import (
	"net/http"
	"strings"
	"sync"
)

type VirtualHosts struct {
	vhosts map[string]http.Handler
	mu     sync.RWMutex
}

func NewVirtualHosts(vhosts map[string]http.Handler) *VirtualHosts {
	v := &VirtualHosts{}
	for hosts, h := range vhosts {
		for _, host := range strings.Split(hosts, " ") {
			if host != "" {
				v.HandleHost(h, host)
			}
		}
	}
	return v
}

func (v *VirtualHosts) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	if v.vhosts != nil {
		if handler, ok := v.vhosts[r.Host]; ok && handler != nil {
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
		for _, h := range strings.Split(host, " ") {
			if h = strings.TrimSpace(h); h != "" {
				v.vhosts[h] = handler
			}
		}
	}
}
