package nvr_format

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"os"
	"strconv"
)
import "time"

//NVR format v1 writer
type NvrWriter struct {
	File   *os.File
	CurPos int64
	Meta   map[string]map[string]int64
	//W io.Writer
}

//create new nvr writer class
func NewWriter() (obj *NvrWriter, err error) {
	obj = &NvrWriter{CurPos: 1048576, Meta: make(map[string]map[string]int64)}
	return
}

//create new nvr file
func (obj *NvrWriter) NewFile(Path string) {
	File, err := os.Create(Path)
	if err != nil {
		return
	}
	obj.File = File
}

//write one frame to file
func (obj *NvrWriter) WriteH264(Data []byte) {
	start := obj.CurPos
	obj.File.Seek(obj.CurPos, 0)
	obj.File.Write([]byte("\000\000\000" + "\r\n"))
	obj.File.Write(Data)
	end, _ := obj.File.Seek(0, 1)
	obj.CurPos = end
	obj.Meta[strconv.FormatInt(time.Now().UnixNano(), 10)] = map[string]int64{"start": start, "end": end, "type": 1}
	obj.WriteMeta()
}

//read file meta data
func (obj *NvrWriter) ReadMeta() {
	//obj.File.Write(Data)

}
func (obj *NvrWriter) WriteMeta() {
	obj.File.Seek(0, 0)
	jsonString, _ := json.Marshal(obj.Meta)
	buf := new(bytes.Buffer)
	gz := gzip.NewWriter(buf)
	gz.Write(jsonString)
	gz.Close()
	//log.Print(buf.Bytes())
	obj.File.Write(buf.Bytes())
}

//close file
func (obj *NvrWriter) Close() {
	obj.File.Close()
}
