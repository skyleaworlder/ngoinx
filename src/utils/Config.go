package utils

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/skyleaworlder/ngoinx/src/config"
	"github.com/tidwall/gjson"
)

// ReadConfig is a tool func to read config file
func ReadConfig(path string) (err error) {
	cfg, err := readConfigFile(path)
	if err != nil {
		log.Println("ngoinx.utils.ReadConfigFile error: readConfigFile failed:", err.Error())
		return err
	}

	initService(cfg)
	return
}

func readConfigFile(path string) (res gjson.Result, err error) {
	fd, err := os.Open(path)
	if err != nil {
		log.Println("ngoinx.utils.readConfigFile error: os.Open failed:", err.Error())
		return gjson.Result{}, err
	}

	contents, err := ioutil.ReadAll(fd)
	if err != nil {
		log.Println("ngoinx.utils.readConfigFile error: ioutil.ReadAll failed:", err.Error())
		return gjson.Result{}, err
	}

	return gjson.Get(string(contents), "service"), nil
}

func initService(cfg gjson.Result) {
	for _, svc := range cfg.Array() {
		s := config.Service{}
		s.Unmarshal(svc)
		config.Svc = append(config.Svc, s)
	}
}
