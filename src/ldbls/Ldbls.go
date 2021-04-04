package ldbls

import (
	"net/http"

	"github.com/skyleaworlder/ngoinx/src/config"
)

// LoadBalancer is an interface
type LoadBalancer interface {
	Init(targets []config.Target) (err error)
	GetAddr(req http.Request) (addr string, err error)
}
