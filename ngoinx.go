package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/skyleaworlder/ngoinx/src/ldbls"
	"github.com/skyleaworlder/ngoinx/src/server"
	"github.com/skyleaworlder/ngoinx/src/utils"
)

func main() {
	svc, _ := utils.ReadConfig("./ngoinx.template.json")
	ldbls.LdblserMapStuffer(svc)

	fd, _ := os.OpenFile("./log/Ngoinx-Server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	logger := utils.LoggerGenerator(&log.TextFormatter{}, fd, log.DebugLevel)
	s := server.NewNgoinxServer(logger, svc)

	if err := s.Serve(); err != nil {
		logger.Fatal("Ngoinx ends its working, err:", err.Error())
	}
}
