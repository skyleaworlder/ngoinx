package ldbls

import (
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/skyleaworlder/ngoinx/src/config"
	"github.com/skyleaworlder/ngoinx/src/utils"
	sm "github.com/umpc/go-sortedmap"
)

// ConsistHash is a struct implement LoadBalancer
type ConsistHash struct {
	Size     int
	No       int
	Compfunc sm.ComparisonFunc
	log      *log.Entry

	HT *sm.SortedMap
}

// NewDefaultConsistHash is to new a default consist hash obj
func NewDefaultConsistHash(size, no int) (c *ConsistHash) {
	compfunc := func(i, j interface{}) bool {
		return i.(*ConHashNode).HID <= j.(*ConHashNode).HID
	}
	logger := log.NewEntry(log.New())
	c = &ConsistHash{Size: size, No: no, Compfunc: compfunc, log: logger, HT: nil}
	return
}

// ConHashNode is a struct
// SN means Serial Number
// HID = sha1(SN + dst)
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
		c.log.Warning(msg)
		return errors.New(msg)
	}

	// each target has its own "dst" and "weight"
	for _, target := range targets {
		// SN is always smaller than "target.Weight"
		for sn := 0; sn < target.Weight; sn++ {
			HID := utils.GetHashID(strconv.Itoa(sn) + target.Dst)
			if err = c.postNode(HID, target.Dst, target.Weight, sn); err != nil {
				return err
			}
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
		c.log.Warning("ngoinx.ldbls.ConsistHash.GetAddr error:", err.Error())
		return "", err
	}
	defer iter.Close()

	var delta uint64 = ^uint64(0)
	for rec := range iter.Records() {
		tmpDelta, err := utils.DeltaUint64(HID, rec.Key.(uint64))
		if err != nil {
			msg := "ngoinx.ldbls.ConsistHash.GetAddr error: conhash.DeltaUint64 failed :"
			c.log.Warning(msg + err.Error())
			return "", errors.New(msg + err.Error())
		}

		// for debug
		c.log.WithFields(log.Fields{"delta": delta, "tmpDelta": tmpDelta}).Info(
			"ConsistHash GetAddr is walking through c.HT and trying to find the suitable Node for request",
		)

		// refresh
		if tmpDelta < delta {
			delta = tmpDelta
			addr = rec.Val.(*ConHashNode).dst
		}
	}
	return addr, nil
}

// SetLogger is to set logger
func (c *ConsistHash) SetLogger(cfg *utils.LoggerConfig) (err error) {
	// e.g LogPath is "./log/", LogFileName is "ConsistHash-1", LogSuffix is ".log"
	// then log file is ./log/ConsistHash-1.log
	logName := cfg.LogPath + cfg.LogFileName + cfg.LogSuffix
	fd, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("ngoinx.ldbls.ConsistHash.SetLogger error: create/open log file", logName, "failed")
		return err
	}
	c.log = utils.LoggerGenerator(cfg.LogFormatter, fd, cfg.LogLevel)
	return nil
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
		c.log.Warning(msg)
		return errors.New(msg)
	}
	return nil
}
