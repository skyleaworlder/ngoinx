package main

import (
	"fmt"
	"testing"

	"github.com/skyleaworlder/ngoinx/src/ldbls"
	"github.com/skyleaworlder/ngoinx/src/server"
	"github.com/skyleaworlder/ngoinx/src/utils"
)

func Test_LdblsMapStuffer(t *testing.T) {
	svc, _ := utils.ReadConfig("../config/ngoinx.template.json")
	ldbls.LdblserMapStuffer(svc)
	fmt.Println(ldbls.LdblserMap)

	server.Serve()
}
