package main

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/skyleaworlder/ngoinx/src/config"
	"github.com/skyleaworlder/ngoinx/src/ldbls"
)

func Test_ConsistHash(t *testing.T) {
	conhash := ldbls.ConsistHash{Size: 8, Compfunc: func(i, j interface{}) bool {
		return i.(*ldbls.ConHashNode).HID <= j.(*ldbls.ConHashNode).HID
	}}

	conhash.Init([]config.Target{
		{Dst: "http://127.0.0.1:10083", Weight: 1},
		{Dst: "http://127.0.0.1:10084", Weight: 2},
	})
	iter, _ := conhash.HT.IterCh()
	for rec := range iter.Records() {
		fmt.Println("rec:", rec.Key, rec.Val)
	}

	url, _ := url.Parse("http://127.0.0.1:10080/api/v1/food")
	req := http.Request{
		Method:        "DELETE",
		URL:           url,
		Header:        http.Header{},
		Body:          http.NoBody,
		ContentLength: 10,
	}
	fmt.Println("request:", req)
	res, _ := conhash.GetAddr(&req)
	fmt.Println("addr:", res)
}
