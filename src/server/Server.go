package server

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/skyleaworlder/ngoinx/src/config"
	"github.com/skyleaworlder/ngoinx/src/ldbls"
	"github.com/skyleaworlder/ngoinx/src/utils"
)

// Server is a struct
type Server struct {
	serverLogger *log.Entry
	svcLoggers   []*log.Entry
}

// NewNgoinxServer is a constructor
// service is always config.Svc
func NewNgoinxServer(logger *log.Entry, service []config.Service) (s *Server) {
	// generate svcLoggers
	// loggers are used in handlerGenerator, in order to get logs when server serving
	svcLoggers := []*log.Entry{}
	for idx, svc := range service {
		logName := "Ngoinx-" + strconv.Itoa(idx) + "-Port-" + strconv.Itoa(int(svc.Listen)) + ".log"
		fd, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			fmt.Println("ngoinx.server.NewNgoinxServer error:", err.Error())
			return nil
		}
		svcLogger := utils.LoggerGenerator(&log.TextFormatter{}, fd, log.DebugLevel)
		svcLoggers = append(svcLoggers, svcLogger)
	}
	return &Server{serverLogger: logger, svcLoggers: svcLoggers}
}

// Serve export for test
func (s *Server) Serve() (err error) {
	s.serverLogger.WithField("status", "INIT_BEGIN").Info("Server begins initializing")
	for idx, svc := range config.Svc {
		handler := handlerGenerator(s.svcLoggers[idx])
		s := &http.Server{
			Addr:           ":" + strconv.Itoa(int(svc.Listen)),
			Handler:        handler,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		go s.ListenAndServe()
	}
	s.serverLogger.WithField("status", "INIT_END").Info("Server finishs initializing, and prepare for serving")
	time.Sleep(1000 * time.Second)
	return errors.New("hahaha")
}

func handlerGenerator(logger *log.Entry) (handler http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Warningln("ngoinx.server.Server.handlerGenerator error: ioutil.ReadAll failed :", err.Error())
			return
		}
		// get scheme://userinfo@host
		addr, err := ldbls.LdblserMap[r.URL.Path].GetAddr(r)
		if err != nil {
			logger.Warningln("ngoinx.server.Server.handlerGenerator error: ldbls.GetAddr failed :", err.Error())
			return
		}

		// cat relayURL
		relayURL := addr + r.URL.Path + r.URL.RawQuery + r.URL.Fragment

		logger.WithFields(log.Fields{
			"url.path":    r.URL.Path,
			"addr chosen": addr,
			"relay URL":   relayURL,
			"map get":     ldbls.LdblserMap[r.URL.Path],
		}).Info("This info appear only if a request comes")

		req, err := http.NewRequest(r.Method, relayURL, strings.NewReader(string(body)))
		logger.Info("request's Header:", req.Header)
		if err != nil {
			logger.Warningln("ngoinx.server.Server.handlerGenerator error: http.NewRequest failed :", err.Error())
			return
		}
		// copy r's Header to req
		req.Header = r.Header

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			logger.Warningln("ngoinx.server.Server.handlerGenerator error: http.DefaultClient.Do failed :", err.Error())
			return
		}
		defer resp.Body.Close()

		// copy resp's Header to rw
		for k, v := range resp.Header {
			rw.Header().Set(k, v[0])
		}
		io.Copy(rw, resp.Body)
	}
}
