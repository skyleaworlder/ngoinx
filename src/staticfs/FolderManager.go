package staticfs

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/skyleaworlder/ngoinx/src/config"
	"github.com/skyleaworlder/ngoinx/src/utils"
)

// SrcPath is defined in ngoinx.json as service.proxy[i].src
type SrcPath string

// FolderManager is a struct to maintain StaticFolders
// each SrcPath has its own *Folder,
// and Folder contains static resource corresponding SrcPath
type FolderManager struct {
	Folders  map[SrcPath]*Folder
	fmLogger *log.Entry
}

// NewDefaultFolderManager is a default constructor
func NewDefaultFolderManager() (sfm *FolderManager) {
	return &FolderManager{Folders: make(map[SrcPath]*Folder), fmLogger: log.NewEntry(log.New())}
}

// Init is to init static folder manager
func (sfm *FolderManager) Init(logger *log.Entry, Service []config.Service) (err error) {
	// assign fmLogger
	sfm.fmLogger = logger

	// generate Folders
	for _, svc := range config.Svc {
		for _, proxy := range svc.Proxies {
			fdr := NewStaticFolder(svc.Static, proxy.Src)
			fdr.mkdir()

			// SetLogger for Folder
			cfg := &utils.LoggerConfig{
				LogPath:      svc.Log,
				LogSuffix:    ".log",
				LogFormatter: &log.TextFormatter{},
				LogLevel:     log.DebugLevel,
				LogFileName:  "StaticFolder-" + strings.ReplaceAll(filepath.Clean(filepath.FromSlash(filepath.Join(svc.Static, proxy.Src))), "\\", "-"),
			}
			cfg.LogOutput, _ = os.OpenFile(cfg.LogPath+cfg.LogFileName+cfg.LogSuffix, os.O_CREATE|os.O_WRONLY, 0755)
			fdr.SetLogger(cfg)
			sfm.Folders[SrcPath(proxy.Src)] = fdr
		}
	}
	return
}

// Serve is a wrapper, only to raise clean
func (sfm *FolderManager) Serve() (err error) {
	return sfm.clean()
}

// StripURLPathPrefix is a tool function
// path: r.URL.Path
// return "" only when (r.URL.Path is src) or fail to find prefix in Folders
// return ok(false) when failing
func (sfm *FolderManager) StripURLPathPrefix(path string) (srcPath SrcPath, dstPath string, ok bool) {
	for srcPath := range sfm.Folders {
		// Folders' keys are all from service[i].proxy[j].src
		// check whether path has the prefix "service[i].proxy[j].src"
		sfm.fmLogger.WithField("r.URL.Path", path).Debugln("this method called in Server.handlerGenerator")
		if dstPath, ok := utils.StripURLPathPrefix(path, string(srcPath)); ok {
			return srcPath, dstPath, ok
		}
	}
	return "", "", false
}

func (sfm *FolderManager) clean() (err error) {
	for running := true; running; {
		sfm.fmLogger.Infoln("ngoinx.staticfs.FolderManager.clean: begin cleaning Folders")
		for _, folder := range sfm.Folders {
			folder.clean()
		}
		time.Sleep(time.Minute)
	}
	return errors.New("ngoinx.staticfs.FolderManager.clean error: clean quit unexpectedly")
}
