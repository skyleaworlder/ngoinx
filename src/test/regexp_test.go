package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/skyleaworlder/ngoinx/src/utils"
)

func Test_regexp_StaticDynamicResourceSeparator(t *testing.T) {
	fmt.Println(utils.StaticSuffixDetermine("/api/v1/haha/1.jpg1", ""))
	fmt.Println(utils.StaticSuffixDetermine("/api/v1/jpg/1.jpg1", ""))
	fmt.Println(utils.StaticSuffixDetermine(".png/v1/haha/1.jpg", ""))
	fmt.Println(utils.StaticSuffixDetermine("/a.htmli/v1/haha/.html.jpg", ""))
	fmt.Println(utils.StaticSuffixDetermine("/.css/.js/.html/.jpg", ""))
	fmt.Println(utils.StaticSuffixDetermine("//v1/h/1.jgp", ""))
	fmt.Println(utils.StaticSuffixDetermine("/api/v1/ha/1.tml", ""))
	//fmt.Println(utils.StaticSuffixDetermine("/api/v1/haha/1.jpg1", ""))
	//fmt.Println(utils.StaticSuffixDetermine("/api/v1/haha/1.jpg1", ""))

	fmt.Println(strings.TrimPrefix("/api/v1/test/hahaha", "/api/v1/test"))
	fmt.Println(strings.TrimPrefix("hahaha", "/api/v1/test"))
}

func Test_regexp_GetStaticFileName(t *testing.T) {
	fmt.Println(utils.GetStaticFileName("./static/v1/api-v1-test/1.js"))
	fmt.Println(utils.GetStaticFileName(".\\static\\v1\\api-v1-test\\1.js"))
}
