package hutils

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func ReadJson(fpath string, i interface{}) error {
	d, e := ioutil.ReadFile(fpath)
	if e != nil {
		return e
	}
	return json.Unmarshal(d, i)
}

func WriteBytes(fpath string, data []byte) error {
	return ioutil.WriteFile(fpath, data, os.ModePerm)
}

func WriteJson(fpath string, i interface{}) error {
	bs, e := json.Marshal(i)
	if e != nil {
		return e
	}
	return WriteBytes(fpath, bs)
}

func AppendBytes(fpath string, bs []byte) error {
	f, e := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if e != nil {
		return e
	}
	defer f.Close()
	_, e = f.Write(bs)
	return e
}
