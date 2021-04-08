package server

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/skyleaworlder/ngoinx/src/config"
	"github.com/skyleaworlder/ngoinx/src/ldbls"
	"github.com/skyleaworlder/ngoinx/src/staticfs"
	"github.com/skyleaworlder/ngoinx/src/utils"
)

// Server is a struct
type Server struct {
	staticServer *staticfs.FolderManager

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

	// new a static folder manager
	sfm := &staticfs.FolderManager{}
	err := sfm.Init(service)
	if err != nil {
		logger.Fatal("ngoinx.server.NewNgoinxServer error: init static folder manager failed")
		return nil
	}

	return &Server{staticServer: sfm, serverLogger: logger, svcLoggers: svcLoggers}
}

// Serve export for test
func (s *Server) Serve() (err error) {
	s.serverLogger.WithField("status", "INIT_BEGIN").Info("Server begins initializing")
	for idx, svc := range config.Svc {
		handler := handlerGenerator(s.svcLoggers[idx], s.staticServer)
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

func handlerGenerator(logger *log.Entry, sfm *staticfs.FolderManager) (handler http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Warningln("ngoinx.server.Server.handlerGenerator error: ioutil.ReadAll failed :", err.Error())
			return
		}

		// Static resources
		// waiting for my implementation
		// utils.StaticSuffixDetermine return isStatic(true) means r.URL.Path suffix satisfy rules
		// e.g. /api/v1/food/1.js or /api/v3/test/.html.css.js
		if isStatic, err := utils.StaticSuffixDetermine(r.URL.Path, ""); isStatic {
			srcPath, result, ok := sfm.StripURLPathPrefix(r.URL.Path)
			// result("") && ok(false) only when proxy doesn't exist, fatal error => return
			// result("") && ok(true) is unknown error (it means r.URL.Path doesn't satisfy static file rules)
			if result == "" && !ok {
				logger.Warningln("ngoinx.server.Server.handlerGenerator error: FolderManager.StripURLPathPrefix failed: prefix do not exist")
				return
			} else if (result == "" && ok) || (result != "" && !ok) {
				logger.Warningln("ngoinx.server.Server.handlerGenerator error: Unknown error about r.URL.Path")
			}
			if err != nil {
				logger.Warningln("ngoinx.server.Server.handlerGenerator error: FolderManager.StripURLPathPrefix failed :", err.Error())
				return
			}

			// check if srcPath in sfm.Folders
			folder, ok := sfm.Folders[srcPath]
			if !ok {
				logger.Warningln("ngoinx.server.Server.handlerGenerator error: srcPath(", srcPath, ") doesn't exist in sfm.Folders")
				return
			}

			// try to open file in folder
			fd, err := folder.Open(filepath.FromSlash(filepath.Clean(result)))
			if err == nil {
				logger.WithFields(log.Fields{"SrcPath": srcPath, "FilePath": filepath.FromSlash(filepath.Clean(result))}).Info(
					"Success Opening FILE, io.Copy will execute later.",
				)
				io.Copy(rw, fd)
				return
			}
			logger.Infoln("ngoinx.server.Server.handlerGenerator error:", result, "do not exist")
			return
		}

		// Dynamic resources
		// get resources from scheme://userinfo@host
		// use r.URL.Path to select Node from MAP maintained by LoadBalancer
		ldblser, ok := ldbls.LdblserMap[r.URL.Path]
		if !ok {
			logger.Warningln("ngoinx.server.Server.handlerGenerator error: r.URL.Path doesn't exist in ldbls.LdblserMap")
			return
		}

		addr, err := ldblser.GetAddr(r)
		if err != nil {
			logger.Warningln("ngoinx.server.Server.handlerGenerator error: ldbls.GetAddr failed :", err.Error())
			return
		}

		// cat relayURL
		relayURL := addr + r.URL.Path + r.URL.RawQuery + r.URL.Fragment

		// for debug
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
