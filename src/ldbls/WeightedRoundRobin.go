package ldbls

import (
	"net/http"

	"github.com/skyleaworlder/ngoinx/src/config"
)

// WeightRoundRobin is a struct implement LoadBalancer
// pointer points at the Node
type WeightRoundRobin struct {
	Size    int
	Nodes   []*RoundRobinNode
	pointer int
}

// NewDefaultWeightedRoundRobin is default constructor
func NewDefaultWeightedRoundRobin(size int) (wr *WeightRoundRobin) {
	return &WeightRoundRobin{Size: size, Nodes: []*RoundRobinNode{}, pointer: 0}
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
func (wr *WeightRoundRobin) Init(targets []config.Target) (err error) {
	for _, target := range targets {
		node := RoundRobinNode{dst: target.Dst, weight: target.Weight, time: 0}
		wr.Nodes = append(wr.Nodes, &node)
	}
	return nil
}

// GetAddr is to implement interface "LoadBalancer"
// ++time == weight => time = 0
// ++pointer == Size => pointer = 0
func (wr *WeightRoundRobin) GetAddr(req *http.Request) (addr string, err error) {
	node := wr.Nodes[wr.pointer]
	addr = node.dst

	// for debug
	// fmt.Println("node.time:", node.time, "wr.pointer:", wr.pointer)
	if node.time++; node.time == node.weight {
		node.time = 0
		if wr.pointer++; wr.pointer == wr.Size {
			wr.pointer = 0
		}
	}
	return addr, nil
}
