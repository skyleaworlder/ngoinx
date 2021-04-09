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
	sfmLogger := log.NewEntry(log.New())

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
	sfm := staticfs.NewDefaultFolderManager()
	err := sfm.Init(sfmLogger, service)
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
		httpServer := &http.Server{
			Addr:           ":" + strconv.Itoa(int(svc.Listen)),
			Handler:        handler,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		go httpServer.ListenAndServe()
	}
	go s.staticServer.Serve()

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

		// e.g. r.URL.Path might be /api/v1/food/1.js
		// 1. get staticResourceName that exists perhaps, if not return ""
		// 2. get srcPath match service[i].proxy[j].src, if not return ""
		// 3. get dstPath via stripPrefix,
		staticResourceName := utils.GetStaticFileName(r.URL.Path)
		srcPath, dstPath, ok := sfm.StripURLPathPrefix(r.URL.Path)

		// Static resources
		// (staticResourceName != "") means r.URL.Path suffix satisfy rules
		// moreover, dstPath & srcPath == "" only when srcPath cannot match
		// e.g. /api/v1/food/1.js or /api/v3/test/.html.css.js
		if staticResourceName != "" {
			// srcPath("") && ok(false) only when proxy doesn't exist, fatal error => return
			if srcPath == "" && !ok {
				logger.WithFields(log.Fields{
					"srcPath":            srcPath,
					"dstPath":            dstPath,
					"r.URL.Path":         r.URL.Path,
					"staticResourcePath": staticResourceName},
				).Warningln("ngoinx.server.Server.handlerGenerator error: FolderManager.StripURLPathPrefix failed: prefix do not exist")
				return
			}

			// dstPath("") is unknown error (it means r.URL.Path doesn't satisfy static file rules)
			if dstPath == "" {
				logger.WithFields(log.Fields{"srcPath": srcPath, "dstPath": dstPath, "ok": ok, "staticResourceName": staticResourceName}).Warningln(
					"ngoinx.server.Server.handlerGenerator error: unknown wrong, dstPath is ''")
				return
			}

			// check if srcPath in sfm.Folders
			folder, ok := sfm.Folders[srcPath]
			if !ok {
				logger.WithFields(log.Fields{"srcPath": srcPath}).Warningln(
					"ngoinx.server.Server.handlerGenerator error: srcPath doesn't exist in sfm.Folders")
				return
			}

			// try to open file in folder
			fd, err := folder.Open(filepath.FromSlash(filepath.Clean(dstPath)))
			if err == nil {
				defer fd.Close()
				logger.WithFields(log.Fields{"SrcPath": srcPath, "FilePath": filepath.FromSlash(filepath.Clean(dstPath))}).Info(
					"Success opening FILE in static resource cache, io.Copy will execute later.")
				io.Copy(rw, fd)
				return
			}

			// file doesn't exist in local storage, log and move to dynamic resources process
			logger.WithFields(log.Fields{"dstPath": dstPath}).Infoln(
				"ngoinx.server.Server.handlerGenerator error: file do not exist, download from remote later")
		}

		// Dynamic resources
		// get resources from scheme://userinfo@host
		// use r.URL.Path to select Node from MAP maintained by LoadBalancer
		ldblser, ok := ldbls.LdblserMap[string(srcPath)]
		if !ok {
			logger.WithFields(log.Fields{"srcPath": srcPath}).Warningln(
				"ngoinx.server.Server.handlerGenerator error: srcPath doesn't exist in ldbls.LdblserMap")
			return
		}

		addr, err := ldblser.GetAddr(r)
		if err != nil {
			logger.Warningln("ngoinx.server.Server.handlerGenerator error: ldbls.GetAddr failed :", err.Error())
			return
		}

		// cat relayURL
		// e.g. service{listen: 10080, proxy:[{src: /api/v1/food, target:[dst: "http://localhost:10081"]}]}
		// http://localhost:10080/api/v1/food => http://localhost:10081
		// http://localhost:10080/api/v1/food/1.js => http://localhost:10081/1.js
		relayURL := addr + dstPath + r.URL.RawQuery + r.URL.Fragment

		// for debug
		logger.WithFields(log.Fields{
			"dstPath":     dstPath,
			"addr chosen": addr,
			"relay URL":   relayURL,
			"map get":     ldbls.LdblserMap[string(srcPath)],
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

		// Static resource also get from resp.Body
		// execute the following block only when static resource doesn't exist
		// (if exists, folder.Open will process)
		if staticResourceName != "" {
			folder, _ := sfm.Folders[srcPath]
			// create file (because it doesn't exist indeed)
			fdc, err := folder.Create(staticResourceName)
			if err != nil {
				logger.Warningln("ngoinx.server.Server.handlerGenerator error: Folder.Create failed :", err.Error())
				return
			}
			io.Copy(fdc, resp.Body)
			fdc.Close()

			// reopen it, in order to use io.Copy(rw, fdo)
			// (since io.Copy(fdc, resp.Body), resp.Body cannot be used anymore)
			fdo, err := folder.Open(staticResourceName)
			if err != nil {
				logger.Warningln("ngoinx.server.Server.handlerGenerator error: Folder.Open failed :", err.Error())
				return
			}
			defer fdo.Close()
			io.Copy(rw, fdo)
			return
		}

		// Dynamic resource process through forrange & io.Copy directly
		// copy resp's Header to rw
		for k, v := range resp.Header {
			rw.Header().Set(k, v[0])
		}
		io.Copy(rw, resp.Body)
	}
}
