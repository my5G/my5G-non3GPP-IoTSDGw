package context

import (
	"github.com/brocaar/lorawan"
	"sync"
	"time"
)


var  (
	APPSKEYTestOnly = lorawan.AES128Key{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	NWKSKEYTestOnly = lorawan.AES128Key{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
)

const (
	Message = "DevPayload"
)

const (
	DefaultTime = "2021-02-17T08:08:30-03:00"
	DefaultTmms = 9223372890
	DefaultTmst = 9223372
	DefaultChan = 0
	DefaultRfch = 1
	DefaultFreq = 916.8
	DefaultStat = 1
	DefaultModu = "LORA"
	DefaultDatr = "SF7BW125"
	DefaultCodr = "4/5"
	DefaultRssi = -57
	DefaultLsnr = 8
)
const(
	FSM_IDLE = iota
	FSM_SEND
	FSM_WAIT
	FSM_RECV
)

// SafeCounter is safe to use concurrently.
type SafeCounter struct {
	mu sync.Mutex
	DevAddrsFlag int
}

type DurationSlice []time.Duration

type Rxpk struct {
	Time string  `json:"time"`
	Tmms int  `json:"tmms"`
	Tmst int  `json:"tmst"`
	Chan uint8   `json:"chan"`
	Rfch uint8   `json:"rfch"`
	Freq float64 `json:"freq"`
	Stat int8     `json:"stat"`
	Modu string  `json:"modu"`
	Datr string  `json:"datr"`
	Codr string  `json:"codr"`
	Rssi int16   `json:"rssi"`
	Lsnr float64  `json:"lsnr"`
	Size uint16   `json:"size"`
	Data string  `json:"data"`
}

type UpStreamJSON struct {
	Rxpk []Rxpk `json:"rxpk"`
}

var idDev uint16
var devAddrsFlag uint16

var counter *SafeCounter
//
// Pega o json avbri o data
/// Converte do base 64 para hex
// Alterar os bytes os bites
func init(){
	idDev = 0
	counter = &SafeCounter{}
}

func (c *SafeCounter) Incremment() (uint16){
	defer c.mu.Unlock()
	c.mu.Lock()
	idDev = idDev + 1
	return idDev
}
func DevIDtoHEx(devId uint16) (lorawan.DevAddr){
	return lorawan.DevAddr{0x00, 0x00, byte(devId >> 8), byte(devId) }
}

func (c *SafeCounter) getAddr() lorawan.DevAddr{
	defer c.mu.Unlock()
	c.mu.Lock()
	c.DevAddrsFlag = c.DevAddrsFlag + 1
	if c.DevAddrsFlag > 9 { c.DevAddrsFlag = 0 }
	data := []lorawan.DevAddr{
		lorawan.DevAddr{0x00,0x00,0x00,0x01}, // DevAddr has 4 bytes
		lorawan.DevAddr{0x00,0x00,0x00,0x02},
		lorawan.DevAddr{0x00,0x00,0x00,0x03},
		lorawan.DevAddr{0x00,0x00,0x00,0x04},
		lorawan.DevAddr{0x00,0x00,0x00,0x05},
		lorawan.DevAddr{0x00,0x00,0x00,0x06},
		lorawan.DevAddr{0x00,0x00,0x00,0x07},
		lorawan.DevAddr{0x00,0x00,0x00,0x08},
		lorawan.DevAddr{0x00,0x00,0x00,0x09},
		lorawan.DevAddr{0x00,0x00,0x00,0x0a},
	}[c.DevAddrsFlag]
	return data
}