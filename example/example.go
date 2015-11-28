package main

import (
	"log"
	"sort"
	"time"

	nvr "./.."
)

func main() {
	////////////////////////////////////new Writer//////////////////////
	objw, _ := nvr.NewWriter()
	//make new data file
	objw.NewFile("test.nvr")
	//get current nano time
	start := time.Now().UnixNano()
	//write frame to data
	//fist
	//0-video 1-audio
	//0-frame 1-key

	objw.WriteH264(0, 1, []byte("hello im test data"))
	objw.WriteH264(1, 1, []byte("hello im big test data!!hello im big test data!!hello im big test data!!hello im big test data!!hello im big test data!!hello im big!!"))
	//zero data
	objw.WriteH264(0, 0, []byte(""))
	objw.WriteH264(0, 1, []byte("d"))
	objw.WriteH264(0, 0, []byte("im strange data(#@$(*@#(*&$*&^@#$*(#$_*!@@)(*#!@)(#*!@#))))"))
	//get current microtime
	end := time.Now().UnixNano()
	//close Writer
	objw.Close()

	///////////////////////////////////////new Reader///////////////////
	objr, _ := nvr.NewReader()
	//open file
	objr.OpenFile("test.nvr")
	packet := objr.ReadTime(start, end)
	for _, v := range sorter(packet) {
		log.Println("time=", v, "type_f=", packet[v]["frame_t"], "type_k=", packet[v]["frame_k"], "payload=", string(packet[v]["payload"]))
	}
	objr.Close()
}
func sorter(data map[string]map[string][]byte) []string {
	mk := make([]string, len(data))
	i := 0
	for k, _ := range data {
		mk[i] = k
		i++
	}
	sort.Strings(mk)
	return mk
}
