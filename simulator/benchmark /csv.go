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
	packet_tx string
	packet_rx string
	recv  string
	timeDuration string

	W *csv.Writer
	F *os.File

	syncWrite sync.Mutex
}

func (h *Metrics) init(){

	h.devId = "DEvID "
	h.packet_tx = "PacketTx"
	h.packet_rx = "PacketRx"
	h.recv =  "Recv"
	h.timeDuration = "Timestamp"

	f, err := os.Create(fmt.Sprintf("LoRaIOTSDGW-Simulator-%s", time.Now().Format("2006.01.02 15:04:05") ))
	if err != nil {
		log.Fatalf("Open Filer to csv writer error to open")
	}

	h.F = f
	h.W = csv.NewWriter(f)
	//defer h.Write.Flush()
	err = h.W.Write([]string{h.devId,h.packet_tx, h.packet_rx, h.timeDuration})
	if err != nil {
		log.Fatalln("error writing header record to file", err)
		return
	}

}

func (h *Metrics) Store(row []string) (error) {
	h.syncWrite.Lock()
	err := h.W.Write(row)
	if err != nil {
		log.Fatalln("error writing  record to file", err)
	}
	h.syncWrite.Unlock()

	return err
}

func (h *Metrics) Close() {
	h.W.Flush()
	err := h.F.Close()
	if err != nil {
		log.Fatalln("error to closed file ", err)
	}
}