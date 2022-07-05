package files

import (
	"bytes"
	"compress/flate"
	"io/ioutil"
)

func Flate(input []byte) []byte {
	var bf = bytes.NewBuffer([]byte{})
	var flater, _ = flate.NewWriter(bf, flate.BestCompression)
	defer flater.Close()
	if _, err := flater.Write(input); err != nil {
		println(err.Error())
		return []byte{}
	}
	if err := flater.Flush(); err != nil {
		println(err.Error())
		return []byte{}
	}
	return bf.Bytes()
}

func UnFlate(input []byte) []byte {
	rdata := bytes.NewReader(input)
	r := flate.NewReader(rdata)
	s, _ := ioutil.ReadAll(r)
	return s
}
