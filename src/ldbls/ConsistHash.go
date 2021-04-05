package ldbls

import (
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"

	"github.com/skyleaworlder/ngoinx/src/config"
	"github.com/skyleaworlder/ngoinx/src/utils"
	sm "github.com/umpc/go-sortedmap"
)

// ConsistHash is a struct implement LoadBalancer
type ConsistHash struct {
	Size     int
	Compfunc sm.ComparisonFunc

	HT *sm.SortedMap
}

// NewDefaultConsistHash is to new a default consist hash obj
func NewDefaultConsistHash(size int) (c *ConsistHash) {
	compfunc := func(i, j interface{}) bool {
		return i.(*ConHashNode).HID <= j.(*ConHashNode).HID
	}
	c = &ConsistHash{Size: size, Compfunc: compfunc, HT: nil}
    return
}

// ConHashNode is a struct
// SN means Serial Number
// HID = sha1(dst + SN)
type ConHashNode struct {
	HID    uint64
	dst    string
	weight int
	SN     int
}

// Init is to implement interface "LoadBalancer"
func (c *ConsistHash) Init(targets []config.Target) (err error) {
	if c.HT = sm.New(c.Size, c.Compfunc); c.HT == nil {
		msg := "ngoinx.ldbls.ConsistHash.Init error: sortedMap.New failed"
		log.Println(msg)
		return errors.New(msg)
	}

	for sn, target := range targets {
		HID := utils.GetHashID(strconv.Itoa(sn) + target.Dst)
		if err = c.postNode(HID, target.Dst, target.Weight, sn); err != nil {
			return err
		}
	}
	return nil
}

// GetAddr is to implement interface "LoadBalancer"
func (c *ConsistHash) GetAddr(req *http.Request) (addr string, err error) {
	buf, _ := httputil.DumpRequest(req, true)
	sha1Sum := sha1.Sum(buf)
	HID := binary.BigEndian.Uint64(sha1Sum[:][:8])
	fmt.Println("hid:", HID)

	iter, err := c.HT.IterCh()
	if err != nil {
		log.Println("ngoinx.ldbls.ConsistHash.GetAddr error:", err.Error())
		return "", err
	}
	defer iter.Close()

	var delta uint64 = ^uint64(0)
	for rec := range iter.Records() {
		tmpDelta, err := utils.DeltaUint64(HID, rec.Key.(uint64))
		if err != nil {
			msg := "ngoinx.ldbls.ConsistHash.GetAddr error: conhash.DeltaUint64 failed :"
			log.Println(msg + err.Error())
			return "", errors.New(msg + err.Error())
		}

		// refresh
		fmt.Print("delta:", delta, "; tmpDelta:", tmpDelta, "\n")
		if tmpDelta < delta {
			// for debug
			// fmt.Print("delta:", tmpDelta, "\n")
			delta = tmpDelta
			addr = rec.Val.(*ConHashNode).dst
		}
	}
	return addr, nil
}

func (c *ConsistHash) postNode(HID uint64, dst string, weight, SN int) (err error) {
	node := &ConHashNode{
		HID:    HID,
		dst:    dst,
		weight: weight,
		SN:     SN,
	}
	if ok := c.HT.Insert(HID, node); !ok {
		msg := "ngoinx.ldbls.ConsistHash.Init warning: HT is full"
		log.Println(msg)
		return errors.New(msg)
	}
	return nil
}
