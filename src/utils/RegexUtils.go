package utils

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// StaticSuffixDetermine is a determiner
func StaticSuffixDetermine(name, rule string) (isStatic bool, err error) {
	if rule == "" {
		rule = "(.html|.css|.js|.jpg|.png|.ico)$"
	}

	reg, _ := regexp.Compile(rule)
	if arr := reg.FindAllString(name, -1); len(arr) > 1 {
		fmt.Println("regexp error! the number of result is:", len(arr))
		return false, errors.New("StaticSuffixDetermine: Unknown Error")
	} else if len(arr) == 0 {
		return false, nil
	}
	return true, nil
}

// StripURLPathPrefix is a tool function
// path: r.URL.Path (e.g. /api/v1/food/1.js, /api/v3/test/favicon.ico)
// prefix: service.proxy[i].src
// this function return result and ok(true) when working well
// while returning result("") and ok(false) when failing
func StripURLPathPrefix(path, prefix string) (result string, ok bool) {
	if strings.HasPrefix(path, prefix) {
		return strings.TrimPrefix(path, prefix), true
	}
	return "", false
}

// GetStaticFileName is a tool function
// if path == "static\v1\api-v1-test\1.js", return 1.js
func GetStaticFileName(path string) (name string) {
	pathFormal := filepath.FromSlash(filepath.Clean(path))
	slashIdx := strings.LastIndex(pathFormal, "\\")
	return pathFormal[slashIdx+1:]
}
