package config

import "github.com/tidwall/gjson"

// Service is a struct (unmarshal from config.json)
type Service struct {
	Listen  uint16
	Log     string
	Static  string
	Proxies []Proxy
}

// Unmarshal is to implement method in interface "Unmarshaler"
func (s *Service) Unmarshal(res gjson.Result) {
	s.Listen = uint16(res.Get("listen").Uint())
	s.Log = res.Get("log").String()
	s.Static = res.Get("static").String()
	for _, proxy := range res.Get("proxy").Array() {
		p := Proxy{}
		p.Unmarshal(proxy)
		s.Proxies = append(s.Proxies, p)
	}
}
