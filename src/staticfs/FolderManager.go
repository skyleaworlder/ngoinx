package staticfs

import (
	"fmt"

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

// StripURLPathPrefix is a tool function
// path: r.URL.Path
// return "" only when (r.URL.Path is src) or fail to find prefix in Folders
// return ok(false) when failing
func (sfm *FolderManager) StripURLPathPrefix(path string) (srcPath SrcPath, result string, ok bool) {
	for srcPath := range sfm.Folders {
		// Folders' keys are all from service[i].proxy[j].src
		// check whether path has the prefix "service[i].proxy[j].src"
		fmt.Println("fordebug, path(r.URL.Path):", path)
		if result, ok := utils.StripURLPathPrefix(path, string(srcPath)); ok {
			return srcPath, result, ok
		}
	}
	return "", "", false
}
