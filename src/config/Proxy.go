package config

import (
	"github.com/tidwall/gjson"
)

// Proxy is a struct (unmarshal from config.file)
type Proxy struct {
	Src    string
	Target []Target
}

// Unmarshal is to implement method in interface "Unmarshaler"
func (p *Proxy) Unmarshal(res gjson.Result) {
	p.Src = res.Get("src").String()
	for _, tgt := range res.Get("target").Array() {
		t := Target{}
		t.Unmarshal(tgt)
		p.Target = append(p.Target, t)
	}
}

// Target is a struct (unmarshal from config.file)
type Target struct {
	Dst    string
	Weight int
}

// Unmarshal is to implement method in interface "Unmarshaler"
func (t *Target) Unmarshal(res gjson.Result) {
	t.Dst = res.Get("dst").String()
	t.Weight = int(res.Get("weight").Int())
}
