package main

import (
	"fmt"
	"testing"

	"github.com/skyleaworlder/ngoinx/src/utils"
)

func Test_Marshal_Unmarshal(t *testing.T) {
	cfg, _ := utils.ReadConfig("../config/ngoinx.template.json")
	svcs := []utils.Service{}
	for _, svc := range cfg.Array() {
		s := utils.Service{}
		fmt.Println("svc:", svc)
		s.Unmarshal(svc)
		svcs = append(svcs, s)
	}
	fmt.Println("svcs:", svcs)
}
