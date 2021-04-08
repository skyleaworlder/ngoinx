package ldbls

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/skyleaworlder/ngoinx/src/config"
	"github.com/skyleaworlder/ngoinx/src/utils"
)

// WeightedRoundRobin is a struct implement LoadBalancer
// pointer points at the Node
type WeightedRoundRobin struct {
	Size    int
	No      int
	Nodes   []*RoundRobinNode
	pointer int
	log     *log.Entry
}

// NewDefaultWeightedRoundRobin is default constructor
func NewDefaultWeightedRoundRobin(size, no int) (wr *WeightedRoundRobin) {
	logger := log.NewEntry(log.New())
	return &WeightedRoundRobin{Size: size, No: no, Nodes: []*RoundRobinNode{}, pointer: 0, log: logger}
}

// RoundRobinNode is a struct
// time means the number of times, which Node has been queried
// time always -le weight
type RoundRobinNode struct {
	dst    string
	weight int
	time   int
}

// Init is to implement interface "LoadBalancer"
func (wr *WeightedRoundRobin) Init(targets []config.Target) (err error) {
	for _, target := range targets {
		node := RoundRobinNode{dst: target.Dst, weight: target.Weight, time: 0}
		wr.Nodes = append(wr.Nodes, &node)
	}
	return nil
}

// GetAddr is to implement interface "LoadBalancer"
// ++time == weight => time = 0
// ++pointer == Size => pointer = 0
func (wr *WeightedRoundRobin) GetAddr(req *http.Request) (addr string, err error) {
	node := wr.Nodes[wr.pointer]
	addr = node.dst

	// for debug
	wr.log.WithFields(log.Fields{"node.time": node.time, "wr.pointer": wr.pointer}).Info(
		"WeightedRoundRobin GetAddr status: node.time and node.weight:", node.weight,
		"wr.pointer and wr.Size:", wr.Size,
	)
	if node.time++; node.time == node.weight {
		node.time = 0
		if wr.pointer++; wr.pointer == wr.Size {
			wr.pointer = 0
		}
	}
	return addr, nil
}

// SetLogger is to implement interface "LoadBalancer"
func (wr *WeightedRoundRobin) SetLogger(cfg *utils.LoggerConfig) (err error) {
	// e.g LogPath is "./log/", LogFileName is "WeightedRoundRobin-1", LogSuffix is ".log"
	// then log file is ./log/WeightedRoundRobin-1.log
	logName := cfg.LogPath + cfg.LogFileName + cfg.LogSuffix
	fd, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("ngoinx.ldbls.WeightedRoundRobin.SetLogger error: create/open log file", logName, "failed")
		return err
	}
	wr.log = utils.LoggerGenerator(cfg.LogFormatter, fd, cfg.LogLevel)
	return
}
