package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/skyleaworlder/ngoinx/src/ldbls"
	"github.com/skyleaworlder/ngoinx/src/server"
	"github.com/skyleaworlder/ngoinx/src/utils"
)

func Test_LdblsMapStuffer(t *testing.T) {
	svc, _ := utils.ReadConfig("../config/ngoinx.template.json")
	ldbls.LdblserMapStuffer(svc)
	fmt.Println(ldbls.LdblserMap)

	fd, _ := os.OpenFile("Ngoinx-Server.log", os.O_CREATE|os.O_WRONLY, 0755)
	logger := utils.LoggerGenerator(&logrus.TextFormatter{}, fd, logrus.DebugLevel)
	s := server.NewNgoinxServer(logger, svc)

	if err := s.Serve(); err != nil {
		fmt.Println("wuhu! end!")
	}
}
