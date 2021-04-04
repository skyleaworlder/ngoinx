package utils

import (
	"github.com/tidwall/gjson"
)

// Proxy is a struct (unmarshal from config.file)
type Proxy struct {
	Src    string
	Target []target
}

// Unmarshal is to implement method in interface "Unmarshaler"
func (p *Proxy) Unmarshal(res gjson.Result) {
	p.Src = res.Get("src").String()
	for _, tgt := range res.Get("target").Array() {
		t := target{}
		t.Unmarshal(tgt)
		p.Target = append(p.Target, t)
	}
}

type target struct {
	Dst    string
	Weight uint8
}

// Unmarshal is to implement method in interface "Unmarshaler"
func (t *target) Unmarshal(res gjson.Result) {
	t.Dst = res.Get("dst").String()
	t.Weight = uint8(res.Get("weight").Uint())
}
