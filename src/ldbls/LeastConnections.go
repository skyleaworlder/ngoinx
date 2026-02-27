package ldbls

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/skyleaworlder/ngoinx/src/config"
	"github.com/skyleaworlder/ngoinx/src/utils"
)

// LeastConnections is a struct implement LoadBalancer
// it selects the node with the lowest served/weight ratio
// ensuring that higher-weight nodes receive proportionally more requests
type LeastConnections struct {
	Size  int
	No    int
	Nodes []*LeastConnNode
	log   *log.Entry
	mu    sync.Mutex
}

// NewDefaultLeastConnections is default constructor
func NewDefaultLeastConnections(size, no int) (lc *LeastConnections) {
	logger := log.NewEntry(log.New())
	return &LeastConnections{Size: size, No: no, Nodes: []*LeastConnNode{}, log: logger}
}

// LeastConnNode is a struct
// served counts the number of requests this node has handled
// weight determines the proportional share of traffic
type LeastConnNode struct {
	dst    string
	weight int
	served int
}

// Init is to implement interface "LoadBalancer"
func (lc *LeastConnections) Init(targets []config.Target) (err error) {
	for _, target := range targets {
		node := LeastConnNode{dst: target.Dst, weight: target.Weight, served: 0}
		lc.Nodes = append(lc.Nodes, &node)
	}
	return nil
}

// GetAddr is to implement interface "LoadBalancer"
// It selects the node with the lowest served/weight ratio.
// When all nodes reach equal ratios, counters are reset to 0.
func (lc *LeastConnections) GetAddr(req *http.Request) (addr string, err error) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	var bestNode *LeastConnNode
	var bestRatio float64 = -1

	for _, node := range lc.Nodes {
		ratio := float64(node.served) / float64(node.weight)

		// for debug
		lc.log.WithFields(log.Fields{"node.dst": node.dst, "node.served": node.served, "node.weight": node.weight, "ratio": ratio}).Info(
			"LeastConnections GetAddr evaluating node",
		)

		if bestNode == nil || ratio < bestRatio {
			bestRatio = ratio
			bestNode = node
		}
	}

	bestNode.served++
	addr = bestNode.dst

	lc.log.WithFields(log.Fields{"chosen": addr, "served": bestNode.served}).Info(
		"LeastConnections GetAddr selected node",
	)

	// check if all nodes have reached equal served/weight ratios
	// if so, reset all counters to allow a fresh cycle
	allEqual := true
	firstRatio := float64(lc.Nodes[0].served) / float64(lc.Nodes[0].weight)
	for _, node := range lc.Nodes[1:] {
		r := float64(node.served) / float64(node.weight)
		if r != firstRatio {
			allEqual = false
			break
		}
	}
	if allEqual {
		for _, node := range lc.Nodes {
			node.served = 0
		}
		lc.log.Info("LeastConnections GetAddr: all nodes reached equal ratio, resetting counters")
	}

	return addr, nil
}

// SetLogger is to implement interface "LoadBalancer"
func (lc *LeastConnections) SetLogger(cfg *utils.LoggerConfig) (err error) {
	// e.g LogPath is "./log/", LogFileName is "LeastConnections-1", LogSuffix is ".log"
	// then log file is ./log/LeastConnections-1.log
	logName := cfg.LogPath + cfg.LogFileName + cfg.LogSuffix
	fd, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("ngoinx.ldbls.LeastConnections.SetLogger error: create/open log file", logName, "failed")
		return err
	}
	lc.log = utils.LoggerGenerator(cfg.LogFormatter, fd, cfg.LogLevel)
	return
}
