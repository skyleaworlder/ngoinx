package staticfs

import (
	"errors"
	"fmt"
	"time"

	"github.com/skyleaworlder/ngoinx/src/config"
	"github.com/skyleaworlder/ngoinx/src/utils"
)

// SrcPath is defined in ngoinx.json as service.proxy[i].src
type SrcPath string

// FolderManager is a struct to maintain StaticFolders
// each SrcPath has its own *Folder,
// and Folder contains static resource corresponding SrcPath
type FolderManager struct {
	Folders map[SrcPath]*Folder
}

// Init is to init static folder manager
func (sfm *FolderManager) Init(Service []config.Service) (err error) {
	sfm.Folders = make(map[SrcPath]*Folder)
	for _, svc := range config.Svc {
		for _, proxy := range svc.Proxies {
			fdr := NewStaticFolder(svc.Static, proxy.Src)
			fdr.mkdir()
			sfm.Folders[SrcPath(proxy.Src)] = fdr
		}
	}
	return
}

// Serve is to raise clean
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
		fmt.Println("fordebug, path(r.URL.Path):", path)
		if dstPath, ok := utils.StripURLPathPrefix(path, string(srcPath)); ok {
			return srcPath, dstPath, ok
		}
	}
	return "", "", false
}

func (sfm *FolderManager) clean() (err error) {
	for running := true; running; {
		fmt.Println("clean!!!!!")
		for _, folder := range sfm.Folders {
			folder.clean()
		}
		time.Sleep(time.Second * 10)
	}
	return errors.New("ngoinx.staticfs.FolderManager.clean quit")
}
