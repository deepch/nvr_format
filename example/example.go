package main

import (
	"log"
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
	objw.WriteH264([]byte("hello im test data"))
	objw.WriteH264([]byte("hello im big test data!!hello im big test data!!hello im big test data!!hello im big test data!!hello im big test data!!hello im big test data!!hello im big test data!!hello im big test data!!hello im big test data!!hello im big test data!!hello im big test data!!hello im big test data!!"))
	//zero data
	objw.WriteH264([]byte(""))
	objw.WriteH264([]byte("d"))
	objw.WriteH264([]byte("im strange data(#@$(*@#(*&$*&^@#$*(#$_*!@@)(*#!@)(#*!@#))))"))
	//get current microtime
	end := time.Now().UnixNano()
	//close Writer
	objw.Close()

	///////////////////////////////////////new Reader///////////////////
	objr, _ := nvr.NewReader()
	//open file
	objr.OpenFile("test.nvr")
	//read file nvr format time to time
	packet := objr.ReadTime(start, end)
	//list of packet
	for k, v := range packet {
		log.Println(k, string(v["payload"]))
	}
	//close Reader
	objr.Close()
}
