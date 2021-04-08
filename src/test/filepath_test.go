package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/skyleaworlder/ngoinx/src/staticfs"
	"github.com/skyleaworlder/ngoinx/src/utils"
)

func Test_filepath(t *testing.T) {
	fmt.Println(os.Args[0])
	config, _ := utils.ReadConfig("../config/ngoinx.template.json")
	sfm := &staticfs.FolderManager{}
	sfm.Init(config)

	for k, v := range sfm.Folders {
		fmt.Println("key:", k, "; val:", v)
	}
}
