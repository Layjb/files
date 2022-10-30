package files

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
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

func XorEncode(bs []byte, keys []byte, cursor int) []byte {
	if len(keys) == 0 {
		return bs
	}

	newbs := make([]byte, len(bs))
	for i, b := range bs {
		newbs[i] = b ^ keys[(i+cursor)%len(keys)]
	}
	return newbs
}

func CreateFile(filename string) (*os.File, error) {
	var err error
	var filehandle *os.File
	if _, err := os.Stat(filename); err == nil { //如果文件存在
		return nil, errors.New("File already exists")
	} else {
		filehandle, err = os.Create(filename) //创建文件
		if err != nil {
			return nil, err
		}
	}
	return filehandle, err
}

func AppendFile(filename string) (*os.File, error) {
	var err error
	var filehandle *os.File
	if _, err := os.Stat(filename); err == nil { //如果文件存在
		filehandle, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return nil, err
		}
	} else {
		filehandle, err = os.Create(filename) //创建文件
		if err != nil {
			return nil, err
		}
	}
	return filehandle, err
}

// Open file from current env and binary path
func Open(filename string) (*os.File, error) {
	f, err := os.Open(filename)
	if err == nil {
		return f, nil
	}

	f, err = os.Open(path.Join(GetExcPath(), filename))
	if err != nil {
		return nil, err
	}
	return f, nil
}

func GetExcPath() string {
	file, _ := exec.LookPath(os.Args[0])
	// 获取包含可执行文件名称的路径
	path, _ := filepath.Abs(file)
	// 获取可执行文件所在目录
	index := strings.LastIndex(path, string(os.PathSeparator))
	ret := path[:index]
	return strings.Replace(ret, "\\", "/", -1) + "/"
}

var Key = []byte{}

func DecryptFile(file *os.File, keys []byte) []byte {
	content, _ := ioutil.ReadAll(file)

	decoded, err := base64.StdEncoding.DecodeString(string(content))
	if err == nil {
		// try to base64 decode, if decode successfully, return data
		return bytes.TrimSpace(decoded)
	}
	// else try to unflate
	decrypted := XorEncode(content, keys, 0)
	if unflated := UnFlate(decrypted); len(unflated) == 0 {
		return bytes.TrimSpace(content)
	} else {
		return bytes.TrimSpace(unflated)
	}
}

func LoadCommonArg(arg string) []byte {
	var content []byte
	f, err := Open(arg)
	if err != nil {
		// if open not found , try base64 decode
		content, err = base64.StdEncoding.DecodeString(arg)
		if err != nil {
			return []byte(arg)
		} else {
			return content
		}
	}
	return DecryptFile(f, Key)
}
