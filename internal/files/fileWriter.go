package files

import (
	"log"
	"os"
)

var WriterChan = make(chan string)
var TotalWrites int

func InitWriter(path string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()

	for {
		s := <-WriterChan
		TotalWrites++
		f.WriteString(s)
	}
}
