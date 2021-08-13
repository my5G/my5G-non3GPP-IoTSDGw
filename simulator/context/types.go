package context

import (
	"encoding/binary"
	"encoding/hex"
	"strconv"
	"strings"
	"time"
)

var idDev uint16

// CompactTime implements time.Time but (un)marshals to and from
// ISO 8601 'compact' format.
type CompactTime time.Time

// MarshalJSON implements the json.Marshaler interface.
func (t CompactTime) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).UTC().Format(`"` + time.RFC3339Nano + `"`)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *CompactTime) UnmarshalJSON(data []byte) error {
	t2, err := time.Parse(`"`+time.RFC3339Nano+`"`, string(data))
	if err != nil {
		return err
	}
	*t = CompactTime(t2)
	return nil
}

// DatR implements the data rate which can be either a string (LoRa identifier)
// or an unsigned integer in case of FSK (bits per second).
type DatR struct {
	LoRa string
	FSK  uint32
}

// MarshalJSON implements the json.Marshaler interface.
func (d DatR) MarshalJSON() ([]byte, error) {
	if d.LoRa != "" {
		return []byte(`"` + d.LoRa + `"`), nil
	}
	return []byte(strconv.FormatUint(uint64(d.FSK), 10)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *DatR) UnmarshalJSON(data []byte) error {
	i, err := strconv.ParseUint(string(data), 10, 32)
	if err != nil {
		d.LoRa = strings.Trim(string(data), `"`)
		return nil
	}
	d.FSK = uint32(i)
	return nil
}

// RXPK contain a RF packet and associated metadata.
type RXPK struct {
	Time CompactTime `json:"time"` // UTC time of pkt RX, us precision, ISO 8601 'compact' format (e.g. 2013-03-31T16:21:17.528002Z)
	Tmst uint32      `json:"tmst"` // Internal timestamp of "RX finished" event (32b unsigned)
	Freq float64     `json:"freq"` // RX central frequency in MHz (unsigned float, Hz precision)
	Chan uint8       `json:"chan"` // Concentrator "IF" channel used for RX (unsigned integer)
	RFCh uint8       `json:"rfch"` // Concentrator "RF chain" used for RX (unsigned integer)
	Stat int8        `json:"stat"` // CRC status: 1 = OK, -1 = fail, 0 = no CRC
	Modu string      `json:"modu"` // Modulation identifier "LORA" or "FSK"
	DatR DatR        `json:"datr"` // LoRa datarate identifier (eg. SF12BW500) || FSK datarate (unsigned, in bits per second)
	CodR string      `json:"codr"` // LoRa ECC coding rate identifier
	RSSI int16       `json:"rssi"` // RSSI in dBm (signed integer, 1 dB precision)
	LSNR float64     `json:"lsnr"` // Lora SNR ratio in dB (signed float, 0.1 dB precision)
	Size uint16      `json:"size"` // RF packet payload size in bytes (unsigned integer)
	Data string      `json:"data"` // Base64 encoded RF packet payload, padded
}


const(
	FSM_IDLE = iota
	FSM_SEND
	FSM_WAIT
	FSM_RECV
)

const (
	PacketType string = "02"
	Payload string = "0000000000000000017b227278706b223a5b7b2274696d65223a22323032312d30322d31375430383a30383a33302d3033" +
		"3a3030222c22746d6d73223a393232333337323839302c22746d7374223a393232333337322c226368616e223a302c2272666368223a312" +
		"c2266726571223a3931362e382c2273746174223a312c226d6f6475223a224c4f5241222c2264617472223a225346374257313235222c22" +
		"636f6472223a22342f35222c2272737369223a2d35372c226c736e72223a372c2273697a65223a31342c2264617461223a2251414541414" +
		"1414141674143347841393846383d227d5d7d"
)

//
// Pega o json avbri o data
/// Converte do base 64 para hex
// Alterar os bytes os bites


func init(){
	idDev = 0
}

func IdEncoding(i uint16) (string){
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, i)
	return hex.EncodeToString(b)
}

func Incremment() (uint16){
	idDev = idDev + 1
	return idDev
}