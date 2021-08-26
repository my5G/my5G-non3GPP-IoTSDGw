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

/*
const (
	PROTOVERSION = "02"

	PacketType string = "02"
	Payload string = "0000000000000000017b227278706b223a5b7b2274696d65223a22323032312d30322d31375430383a30383a33302d3033" +
		"3a3030222c22746d6d73223a393232333337323839302c22746d7374223a393232333337322c226368616e223a302c2272666368223a312" +
		"c2266726571223a3931362e382c2273746174223a312c226d6f6475223a224c4f5241222c2264617472223a225346374257313235222c22" +
		"636f6472223a22342f35222c2272737369223a2d35372c226c736e72223a372c2273697a65223a31342c2264617461223a2251414541414" +
		"1414141674143347841393846383d227d5d7d"
)
*/

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
	Rxpk []Rxpk`json:"rxpk"`
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

/*
func IdEncoding(i uint16) (string){
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, i)
	return hex.EncodeToString(b)
}
*/

func Incremment() (uint16){
	idDev = idDev + 1
	return idDev
}

func (c * SafeCounter) getAddr() lorawan.DevAddr{
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