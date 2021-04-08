package staticfs

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/skyleaworlder/ngoinx/src/utils"
)

// Folder is a struct implement http.FileSystem
// path is folder's path, if folder is /usr/local/haha, "path" equal /usr/local/haha
// NOTICE: "path" is real path in computer-fs
// path is defined as service.static + strings.ReplaceAll(service.proxy[i].src[1:], "/", "-")
type Folder struct {
	path string
	log  *log.Entry
}

// NewDefaultFolder is a default constructor
func NewDefaultFolder(path string) (f *Folder) {
	return &Folder{path: path, log: log.NewEntry(log.New())}
}

// NewStaticFolder is a constructor
// path <= filepath.Join(svc.Static, strings.ReplaceAll(proxy.Src[1:], "/", "-"))
// log also use default value (log.NewEntry(log.New()))
func NewStaticFolder(staticPath, srcPath string) (sf *Folder) {
	return &Folder{
		path: filepath.Join(staticPath, strings.ReplaceAll(srcPath[1:], "/", "-")),
		log:  log.NewEntry(log.New()),
	}
}

// Open is to implement interface "http.FileSystem"
// name should use "\", instead of "/", but filepath.FromSlash can also solve this unstandard input
// e.g. name is "list.json", and s.path is ".\static\v1\api-v1-food"
// then fileName is ".\static\v1\api-v1-food\list.json"
func (sf *Folder) Open(name string) (file http.File, err error) {
	fileName := filepath.Join(sf.path, filepath.FromSlash(filepath.Clean(name)))
	file, err = os.Open(fileName)
	if err != nil {
		sf.log.WithField("fileName", fileName).Warningln(
			"ngoinx.staticfs.Folder.Open error: os.Open failed to open file, return err:", err.Error())
	}
	return
}

// Create is a method used to create file in static folder (managed by Folder struct)
func (sf *Folder) Create(name string) (file *os.File, err error) {
	fileName := filepath.Join(sf.path, filepath.FromSlash(filepath.Clean(name)))
	file, err = os.Create(fileName)
	if err != nil {
		sf.log.WithField("fileName", fileName).Warningln(
			"ngoinx.staticfs.Folder.Create error: os.Create failed to create file, return err:", err.Error())
	}
	return
}

// SetLogger is to implement interface "Loggerable"
func (sf *Folder) SetLogger(cfg *utils.LoggerConfig) (err error) {
	// e.g LogPath is "./log/", LogFileName is "StaticFolder-v1-api-v1-food", LogSuffix is ".log"
	// then log file is ./log/StaticFolder-v1-api-v1-food.log
	logName := cfg.LogPath + cfg.LogFileName + cfg.LogSuffix
	fd, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		sf.log.WithField("logName", logName).Warningln(
			"ngoinx.staticfs.Folder.SetLogger error: create/open log file failed")
		return err
	}
	sf.log = utils.LoggerGenerator(cfg.LogFormatter, fd, cfg.LogLevel)
	return nil
}

// mkdir wrap os.MkdirAll function
// mode default 0755
func (sf *Folder) mkdir() (err error) {
	err = os.MkdirAll(sf.path, 0755)
	if err != nil {
		sf.log.WithField("sf.path", sf.path).Fatalln(
			"ngoinx.staticfs.Folder.mkdir Fatal error: making directory for Static Folder failed, err:", err.Error())
	}
	return
}

// clean is to delete some files that are not to standard
func (sf *Folder) clean() {
	filepath.Walk(sf.path, func(path string, info os.FileInfo, err error) error {
		isDir := info.IsDir()
		bodyNull := (info.Size() == 0)
		outDated := (time.Now().Sub(info.ModTime()) >= time.Minute)
		if isDelete := !isDir && (bodyNull || outDated); isDelete {
			err = os.Remove(path)
			if err != nil {
				sf.log.WithField("info.Name", info.Name()).Warningln(
					"ngoinx.staticfs.Folder.clean error: remove file failed, return err:", err.Error())
				return err
			}
			sf.log.WithField("info.Name", info.Name()).Warningln("ngoinx.staticfs.Folder.clean: remove file SUCCESS")
		}
		return nil
	})
}
