package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/skyleaworlder/ngoinx/src/staticfs"
	"github.com/skyleaworlder/ngoinx/src/utils"
)

func Test_filepath(t *testing.T) {
	fmt.Println(os.Args[0])
	config, _ := utils.ReadConfig("../config/ngoinx.template.json")
	sfm := &staticfs.FolderManager{}
	sfmLogger := logrus.NewEntry(logrus.New())
	sfm.Init(sfmLogger, config)

	for k, v := range sfm.Folders {
		fmt.Println("key:", k, "; val:", v)
	}
}

func Test_Walk(t *testing.T) {
	filepath.Walk("./static", func(path string, info os.FileInfo, err error) error {
		fmt.Println(path, info)
		return nil
	})
}
