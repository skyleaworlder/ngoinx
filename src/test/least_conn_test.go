package main

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/skyleaworlder/ngoinx/src/config"
	"github.com/skyleaworlder/ngoinx/src/ldbls"
)

func Test_LeastConnections_EqualWeight(t *testing.T) {
	lc := ldbls.NewDefaultLeastConnections(3, 0)
	lc.Init([]config.Target{
		{Dst: "http://127.0.0.1:10081", Weight: 1},
		{Dst: "http://127.0.0.1:10082", Weight: 1},
		{Dst: "http://127.0.0.1:10083", Weight: 1},
	})

	url, _ := url.Parse("http://127.0.0.1:10080/api/v1/test")
	req := &http.Request{
		Method: "GET",
		URL:    url,
		Header: http.Header{},
		Body:   http.NoBody,
	}

	// With equal weights, each node should be selected once before any repeats
	counts := map[string]int{}
	for i := 0; i < 9; i++ {
		addr, err := lc.GetAddr(req)
		if err != nil {
			t.Fatal("GetAddr failed:", err)
		}
		counts[addr]++
	}

	// Each node should have been selected exactly 3 times
	for dst, count := range counts {
		if count != 3 {
			t.Errorf("Expected 3 requests for %s, got %d", dst, count)
		}
	}
	fmt.Println("EqualWeight counts:", counts)
}

func Test_LeastConnections_WeightedDistribution(t *testing.T) {
	lc := ldbls.NewDefaultLeastConnections(3, 0)
	lc.Init([]config.Target{
		{Dst: "http://127.0.0.1:10081", Weight: 1},
		{Dst: "http://127.0.0.1:10082", Weight: 2},
		{Dst: "http://127.0.0.1:10083", Weight: 3},
	})

	url, _ := url.Parse("http://127.0.0.1:10080/api/v1/test")
	req := &http.Request{
		Method: "GET",
		URL:    url,
		Header: http.Header{},
		Body:   http.NoBody,
	}

	// Run through enough requests to see the weight distribution
	// Weights 1:2:3 means in one full cycle:
	// node1 gets 1, node2 gets 2, node3 gets 3 = 6 total
	counts := map[string]int{}
	for i := 0; i < 6; i++ {
		addr, err := lc.GetAddr(req)
		if err != nil {
			t.Fatal("GetAddr failed:", err)
		}
		counts[addr]++
	}

	fmt.Println("WeightedDistribution counts:", counts)

	// Verify proportional distribution
	if counts["http://127.0.0.1:10081"] != 1 {
		t.Errorf("Expected 1 request for :10081 (weight 1), got %d", counts["http://127.0.0.1:10081"])
	}
	if counts["http://127.0.0.1:10082"] != 2 {
		t.Errorf("Expected 2 requests for :10082 (weight 2), got %d", counts["http://127.0.0.1:10082"])
	}
	if counts["http://127.0.0.1:10083"] != 3 {
		t.Errorf("Expected 3 requests for :10083 (weight 3), got %d", counts["http://127.0.0.1:10083"])
	}
}

func Test_LeastConnections_SingleTarget(t *testing.T) {
	lc := ldbls.NewDefaultLeastConnections(1, 0)
	lc.Init([]config.Target{
		{Dst: "http://127.0.0.1:10081", Weight: 1},
	})

	url, _ := url.Parse("http://127.0.0.1:10080/api/v1/test")
	req := &http.Request{
		Method: "GET",
		URL:    url,
		Header: http.Header{},
		Body:   http.NoBody,
	}

	// Single target should always return the same address
	for i := 0; i < 5; i++ {
		addr, err := lc.GetAddr(req)
		if err != nil {
			t.Fatal("GetAddr failed:", err)
		}
		if addr != "http://127.0.0.1:10081" {
			t.Errorf("Expected http://127.0.0.1:10081, got %s", addr)
		}
	}
	fmt.Println("SingleTarget: all requests correctly routed to single node")
}
