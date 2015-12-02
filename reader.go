package nvr_format

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
)

//NVR format v1 writer
type NvrReader struct {
	File *os.File
	Meta map[string]map[string]int64
	//W io.Writer
}

//create new nvr writer class
func NewReader() (obj *NvrReader, err error) {
	obj = &NvrReader{}
	return
}

//create new nvr file
func (obj *NvrReader) OpenFile(Path string) {
	File, err := os.Open(Path)
	if err != nil {
		panic("not file")
	}
	obj.File = File
	obj.ReadMeta()
}

//read file meta data
func (obj *NvrReader) ReadMeta() {
	header := make([]byte, 1048576)
	obj.File.Read(header)
	gz, _ := gzip.NewReader(bytes.NewBuffer(header))
	plaintext, _ := ioutil.ReadAll(gz)
	var y map[string]map[string]int64
	json.Unmarshal(plaintext, &y)
	obj.Meta = y
}

func (obj *NvrReader) ReadTime(start, end int64) map[string]map[string][]byte {
	y := make(map[string]map[string][]byte)
	for k, v := range obj.Meta {
		i, _ := strconv.ParseInt(k, 10, 64)
		if i >= start && i <= end {
			diff := v["e"] - v["s"]
			header := make([]byte, diff)
			obj.File.Seek(v["s"], 0)
			obj.File.Read(header)
			y[k] = make(map[string][]byte)
			y[k]["payload"] = header
			y[k]["t"] = []byte(string(v["t"]))
			y[k]["k"] = []byte(string(v["k"]))
		}
	}

	return y
}

//close file
func (obj *NvrReader) Close() {
	obj.File.Close()
}
