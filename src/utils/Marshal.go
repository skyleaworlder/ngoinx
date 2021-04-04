package utils

import "github.com/tidwall/gjson"

// MarshalUnmarshaler is balabala
type MarshalUnmarshaler interface {
	Marshaler
	Unmarshaler
}

// Marshaler is an interface
// turn v(interface{}) => gjson.Result
type Marshaler interface {
	Marshal(v interface{}) (res gjson.Result)
}

// Unmarshaler is an interface
//
// func (v interface{}) Unmarshal(res gjson.Result)
// turn gjson.Result => v(interface{})
type Unmarshaler interface {
	Unmarshal(res gjson.Result)
}
