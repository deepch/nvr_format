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
		return
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

func (obj *NvrReader) ReadTime(start, end int64) map[int64]map[string][]byte {
	y := make(map[int64]map[string][]byte)
	for k, v := range obj.Meta {
		i, _ := strconv.ParseInt(k, 10, 64)
		if i >= start && i <= end {
			diff := v["end"] - v["start"]
			header := make([]byte, diff)
			obj.File.Seek(v["start"], 0)
			obj.File.Read(header)
			y[i] = make(map[string][]byte)
			y[i]["payload"] = header[5:]
			y[i]["type"] = []byte(string(v["type"]))
		}
	}
	return y
}

//close file
func (obj *NvrReader) Close() {
	obj.File.Close()
}
