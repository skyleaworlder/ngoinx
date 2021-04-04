package main

import (
	"fmt"
	"testing"

	"github.com/skyleaworlder/ngoinx/src/config"
	"github.com/skyleaworlder/ngoinx/src/utils"
)

func Test_Marshal_Unmarshal(t *testing.T) {
	utils.ReadConfig("../config/ngoinx.template.json")
	fmt.Println("svcs:", config.Svc)
}
