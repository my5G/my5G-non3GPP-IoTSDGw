package main

import (
	"encoding/csv"
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"sync"
	"time"
)

const ROOT = "/metrics"

type Metrics struct {

	msgType string
	packet_seq int
	byteSize int
	time time.Time

	W *csv.Writer
	F *os.File

	syncWrite sync.Mutex
}

func (h *Metrics) Init(){

	id := uuid.New()
	f, err := os.Create(fmt.Sprintf(ROOT+"/IOTSDGW-Forwarder-%s.csv", id.String() ))
	if err != nil {
		log.Fatalf("Open Filer to csv writer error to open")
	}

	h.F = f
	h.W = csv.NewWriter(f)
	//defer h.Write.Flush()

	err = h.W.Write( []string{ "type", "Seq", "size", "Timestamp" })
	if err != nil {
		log.Fatalln("error writing header record to file", err)
		return
	}
}

func StoreInfo( msgType string, packet_seq int, byteSize int, time float64  ) []string {
	return []string{
		fmt.Sprintf("%s", msgType),
		fmt.Sprintf("%d", packet_seq),
		fmt.Sprintf("%d", byteSize),
		fmt.Sprintf("%f", time),
	}
}

func (h *Metrics) Store(row []string) (error) {
	defer h.syncWrite.Unlock()
	h.syncWrite.Lock()
	err := h.W.Write(row)
	if err != nil {
		log.Fatalln("error writing  record to file", err)
	}

	return err
}

func (h *Metrics) Close() {

	h.W.Flush()

	err := h.F.Close()
	if err != nil {
		log.Fatalln("error to closed file ", err)
	}
}