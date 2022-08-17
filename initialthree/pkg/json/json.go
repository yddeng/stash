package json

import (
	jsoniter "github.com/json-iterator/go"
)

var Marshal func(interface{}) ([]byte, error) = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal
var Unmarshal func([]byte, interface{}) error = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal
