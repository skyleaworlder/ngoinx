package staticfs

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Folder is a struct implement http.FileSystem
// path is folder's path, if folder is /usr/local/haha, "path" equal /usr/local/haha
// NOTICE: "path" is real path in computer-fs
// path is defined as service.static + strings.ReplaceAll(service.proxy[i].src[1:], "/", "-")
type Folder struct {
	path string
}

// NewDefaultFolder is a default constructor
func NewDefaultFolder(path string) (f *Folder) {
	return &Folder{path: path}
}

// NewStaticFolder is a constructor
// path <= filepath.Join(svc.Static, strings.ReplaceAll(proxy.Src[1:], "/", "-"))
func NewStaticFolder(staticPath, srcPath string) (sf *Folder) {
	return &Folder{path: filepath.Join(staticPath, strings.ReplaceAll(srcPath[1:], "/", "-"))}
}

// Open is to implement interface "http.FileSystem"
// name should use "\", instead of "/", but filepath.FromSlash can also solve this unstandard input
// e.g. name is "list.json", and s.path is ".\static\v1\api-v1-food"
// then fileName is ".\static\v1\api-v1-food\list.json"
func (sf *Folder) Open(name string) (file http.File, err error) {
	fileName := filepath.Join(sf.path, filepath.FromSlash(filepath.Clean(name)))
	file, err = os.Open(fileName)
	return
}

// mkdir wrap os.MkdirAll function
// mode default 0755
func (sf *Folder) mkdir() (err error) {
	return os.MkdirAll(sf.path, 0755)
}
