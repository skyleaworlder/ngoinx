package ldbls

import (
	"net/http"

	"github.com/skyleaworlder/ngoinx/src/config"
)

var (
	// LdblserMap is a global map
	LdblserMap map[string]LoadBalancer = make(map[string]LoadBalancer)
)

// LoadBalancer is an interface
type LoadBalancer interface {
	Init(targets []config.Target) (err error)
	GetAddr(req *http.Request) (addr string, err error)
}

// LdblserMapStuffer is a tool func to fill ldblsermap
func LdblserMapStuffer() {
	for _, svc := range config.Svc {
		for _, proxy := range svc.Proxies {
			var ldblser LoadBalancer
			if len(proxy.Target) >= 4 {
				ldblser = NewDefaultConsistHash(len(proxy.Target))
			} else {
				ldblser = &WeightRoundRobin{Size: len(proxy.Target), Nodes: []*RoundRobinNode{}}
			}
			ldblser.Init(proxy.Target)
			LdblserMap[proxy.Src] = ldblser
		}
	}
}

// LdblserGenerator will generate *LoadBalancer
func LdblserGenerator(ldblser *LoadBalancer) {
	return
}
