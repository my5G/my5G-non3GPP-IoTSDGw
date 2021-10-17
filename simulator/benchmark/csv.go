package benchmark

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
	devId string
	msgType string
	packet_seq string
	recv  string
	time string

	totalPacket int
	totalTime float64

	W *csv.Writer
	F *os.File

	W2 *csv.Writer
	F2 *os.File

	syncWrite sync.Mutex
	syncWriteResume sync.Mutex
}

func (h *Metrics) Init(){

	id := uuid.New()


	h.devId = "DevID "
	h.msgType  = "type"
	h.packet_seq = "Seq"
	h.recv = "recv"
	h.time = "Timestamp"

	f, err := os.Create(ROOT + fmt.Sprintf("/LoRaIOTSDGW-Simulator-%s.csv", id.String() ))
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


	f2, err := os.Create(ROOT + fmt.Sprintf("/LoRaIOTSDGW-Sim-Resume%s.csv", time.Now().Format("2006.01.02 15:04:05") ))
	if err != nil {
		log.Fatalf("Open Filer to csv writer error to open")
	}

	h.F2 = f2
	h.W2 = csv.NewWriter(f2)
	//defer h.Write.Flush()

	err = h.W2.Write([]string{h.devId, h.msgType, "TotalPacket","TotalTime"})
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

func (h *Metrics) StoreResume(row []string) (error) {
	defer h.syncWriteResume.Unlock()
	h.syncWriteResume.Lock()
	err := h.W2.Write(row)
	if err != nil {
		log.Fatalln("error writing  record to file", err)
	}

	return err
}

func (h *Metrics) Close() {
	h.W.Flush()
	h.W2.Flush()
	err := h.F.Close()
	if err != nil {
		log.Fatalln("error to closed file ", err)
	}
	err = h.F2.Close()
	if err != nil {
		log.Fatalln("error to closed file ", err)
	}
}