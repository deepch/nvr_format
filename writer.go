package nvr_format

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)
import "time"

//NVR format v1 writer
type NvrWriter struct {
	File   *os.File
	CurPos int64
	Meta   map[string]map[string]int64
	Count  int64
	//W io.Writer
}

type DataNvrMap struct {
	s uint64
	e uint64
	t byte
	k byte
	//W io.Writer
}

//create new nvr writer class
func NewWriter() (obj *NvrWriter, err error) {
	obj = &NvrWriter{CurPos: 1048576, Meta: make(map[string]map[string]int64)}
	return
}

//create new nvr file
func (obj *NvrWriter) NewFile(Path string) {

	if _, err := os.Stat(Path); os.IsNotExist(err) {
		file, _ := os.Create(Path)
		obj.File = file
	} else {
		file, _ := os.OpenFile(Path, os.O_RDWR, 0777)
		obj.File = file
		obj.ReadMeta()
		end, _ := obj.File.Seek(0, 2)
		obj.CurPos = end
	}
}

//write one frame to file
func (obj *NvrWriter) WriteH264(frame_t int64, frame_k int64, Data []byte) {
	obj.Count++
	start := obj.CurPos
	obj.File.Seek(obj.CurPos, 0)
	obj.File.Write(Data)
	end, _ := obj.File.Seek(0, 1)
	obj.CurPos = end
	obj.Meta[strconv.FormatInt(time.Now().UnixNano(), 10)] = map[string]int64{"s": start, "e": end, "t": frame_t, "k": frame_k}
	if obj.Count%500 == 0 {
		log.Print("write meta")
		obj.WriteMeta()
	}
}

//read file meta data
func (obj *NvrWriter) ReadMeta() {
	header := make([]byte, 1048576)
	obj.File.Read(header)
	gz, err := gzip.NewReader(bytes.NewBuffer(header))
	if err != nil {
		return
	}
	plaintext, _ := ioutil.ReadAll(gz)
	var y map[string]map[string]int64
	json.Unmarshal(plaintext, &y)
	obj.Meta = y
}
func (obj *NvrWriter) WriteMeta() {
	obj.File.Seek(0, 0)
	jsonString, _ := json.Marshal(obj.Meta)
	buf := new(bytes.Buffer)
	gz, _ := gzip.NewWriterLevel(buf, gzip.BestSpeed)
	gz.Write(jsonString)
	gz.Close()
	obj.File.Write(buf.Bytes())
	return
}

//close file
func (obj *NvrWriter) Close() {
	obj.WriteMeta()
	obj.File.Close()
}
