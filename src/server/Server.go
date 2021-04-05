package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/skyleaworlder/ngoinx/src/config"
	"github.com/skyleaworlder/ngoinx/src/ldbls"
)

// Serve export for test
func Serve() {
	for _, svc := range config.Svc {
		handler := handlerGenerator()
		s := &http.Server{
			Addr:           ":" + strconv.Itoa(int(svc.Listen)),
			Handler:        handler,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		go s.ListenAndServe()
	}
	fmt.Println("serve finish init.")
	time.Sleep(1000 * time.Second)
}

func handlerGenerator() (handler http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("ngoinx.server.Server.handlerGenerator error: ioutil.ReadAll failed :", err.Error())
			return
		}
		// get scheme://userinfo@host
		fmt.Println("url.path", r.URL.Path)
		fmt.Println("map get:", ldbls.LdblserMap, ldbls.LdblserMap[r.URL.Path])
		addr, err := ldbls.LdblserMap[r.URL.Path].GetAddr(r)
		fmt.Println("addr chosen:", addr)
		if err != nil {
			log.Println("ngoinx.server.Server.handlerGenerator error: ldbls.GetAddr failed :", err.Error())
			return
		}

		// cat relayURL
		relayURL := addr + "/" + r.URL.Path + r.URL.RawQuery + r.URL.Fragment
		req, err := http.NewRequest(r.Method, relayURL, strings.NewReader(string(body)))
		fmt.Println(req.Header)
		if err != nil {
			log.Println("ngoinx.server.Server.handlerGenerator error: http.NewRequest failed :", err.Error())
			return
		}
		// copy r's Header to req
		req.Header = r.Header

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("ngoinx.server.Server.handlerGenerator error: http.DefaultClient.Do failed :", err.Error())
			return
		}
		defer resp.Body.Close()

		// copy resp's Header to rw
		for k, v := range resp.Header {
			rw.Header().Set(k, v[0])
		}
		io.Copy(rw, resp.Body)
	}
}
