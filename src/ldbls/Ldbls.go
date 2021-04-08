package ldbls

import (
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/skyleaworlder/ngoinx/src/config"
	"github.com/skyleaworlder/ngoinx/src/utils"
)

var (
	// LdblserMap is a global map
	LdblserMap map[string]LoadBalancer = make(map[string]LoadBalancer)
)

// LoadBalancer is an interface
type LoadBalancer interface {
	// Init is a init method
	// every struct implements LoadBalancer might own NewDefaultLoadBalancer as default constructor
	// after getting LoadBalancer from NewDefaultLoadBalancer
	// (LoadBalancer).Init is also necessary
	// NewDefaultLoadBalancer might need several parameters,
	// but it should only pass each parameter to property
	// Init will do with much more complicated work
	Init(targets []config.Target) (err error)

	// GetAddr is a method to Get URL.Scheme+"://"+URL.Host
	// inner data structure has been initialized well in Init method
	// GetAddr only pass request as parameter, and then return addr
	// in server package,
	// use "ldbls.LdblserMap[req.URL.Path].GetAddr(r)" to get addr
	GetAddr(req *http.Request) (addr string, err error)

	// Loggerable need to implement method "SetLogger"
	utils.Loggerable
}

// LdblserMapStuffer is a tool func to fill ldblsermap
// len(target) \in [1, 3] => use WeightedRoundRobin
// len(target) \in [4, \inf] => use ConsistHash
func LdblserMapStuffer(service []config.Service) {
	for _, svc := range service {
		for no, proxy := range svc.Proxies {
			var ldblser LoadBalancer
			cfg := &utils.LoggerConfig{LogPath: svc.Log, LogSuffix: ".log", LogFormatter: &log.TextFormatter{}, LogLevel: log.DebugLevel}
			if len(proxy.Target) >= 4 {
				cfg.LogFileName = "ConsistHash-" + strconv.Itoa(no)
				cfg.LogOutput, _ = os.OpenFile(cfg.LogPath+cfg.LogFileName+cfg.LogSuffix, os.O_CREATE|os.O_WRONLY, 0755)
				ldblser = NewDefaultConsistHash(len(proxy.Target), no)
				ldblser.SetLogger(cfg)
			} else {
				cfg.LogFileName = "WeightedRoundRobin-" + strconv.Itoa(no)
				cfg.LogOutput, _ = os.OpenFile(cfg.LogPath+cfg.LogFileName+cfg.LogSuffix, os.O_CREATE|os.O_WRONLY, 0755)
				ldblser = NewDefaultWeightedRoundRobin(len(proxy.Target), no)
				ldblser.SetLogger(cfg)
			}
			ldblser.Init(proxy.Target)
			LdblserMap[proxy.Src] = ldblser
		}
	}
}
