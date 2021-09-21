package benchmark

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Metrics struct {
	devId string
	msgType string
	packet_seq string
	recv  string
	time string

	W *csv.Writer
	F *os.File

	syncWrite sync.Mutex
}

func (h *Metrics) Init(){

	h.devId = "DevID "
	h.msgType  = "type"
	h.packet_seq = "Seq"
	h.recv = "recv"
	h.time = "Timestamp"

	f, err := os.Create(fmt.Sprintf("LoRaIOTSDGW-Simulator-%s.csv", time.Now().Format("2006.01.02 15:04:05") ))
	if err != nil {
		log.Fatalf("Open Filer to csv writer error to open")
	}

	h.F = f
	h.W = csv.NewWriter(f)
	//defer h.Write.Flush()

	err = h.W.Write([]string{h.devId, h.msgType, h.packet_seq, h.recv, h.time})
	if err != nil {
		log.Fatalln("error writing header record to file", err)
		return
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