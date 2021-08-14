package context

import (
	"encoding/hex"
	"fmt"
	"github.com/brocaar/lorawan"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"sync"
	"time"
)

var ctx DevContext

type Device struct {
	Id uint16
	Packet_tx, Packet_rx uint8
	Timeout time.Timer
	Latency time.Time
	PacketLoss int
	FsmState int
	confirmed bool

	// AppKey.
	//appKey lorawan.AES128Key
	// Application session-key.
	appSKey lorawan.AES128Key
	// Network session-key.
	nwkSKey lorawan.AES128Key
	// devAddr
	devAddr lorawan.DevAddr
	// DevEUI.
	//devEUI lorawan.EUI64
	// FPort used for sending uplinks.
	fPort uint8
	// Uplink frame-counter.
	fCntUp uint32
	// Downlink frame-counter.
	fCntDown uint32
	// Payload (plaintext) which the device sends as uplink.
	payload []byte
}

type DevContext struct{
	DevicesPool  sync.Map
	ForwarderConn *net.UDPAddr
}

func DevicesContext_Self() *DevContext{
	return &ctx
}

func (device *Device) init( id uint16){
	device.Id = id
}

func (device *Device) MakeNewPayload() ([]byte, bool){
	hexMessage := fmt.Sprintf("%s%s%s",
		PacketType,
		IdEncoding(device.Id),
		Payload )

	msg, err := hex.DecodeString(hexMessage)
	if err != nil {
		return nil, false
	}
	return msg, true
}

func ( ctx *DevContext) NewDevice() *Device {
	valueID := Incremment()
	if  valueID < 0  {
		log.Fatalf("Dev Id code Error")
	}
	device := new(Device)
	device.init(valueID)
	ctx.DevicesPool.Store(valueID, device)
	device.FsmState = FSM_IDLE
	return device
}

func ( ctx *DevContext) DeviceLoad(id uint16) (*Device, bool) {
	device, ok := ctx.DevicesPool.Load(id)
	if ok {
		return device.(*Device), ok
	} else {
		return nil, ok
	}
}

func CreateDevicesForSimulate(devicesLen int){
	if devicesLen < 1 {
		log.Fatalf("Number of devices is not valid")
	}

	for i := 0; i < devicesLen; i++ {
		DevicesContext_Self().NewDevice()
	}
}

func ( ctx *DevContext) ConfigSocketUDPAddr( ipAddr string, port int)(bool){
	serverAddr,err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d",
		ipAddr, port))
	if err != nil {
		log.Fatalf("Config Bind IPAddr Error %v", err)
		return false
	}
	ctx.ForwarderConn = serverAddr

	return true
}

// dataUp sends an data uplink.
func (d *Device) dataUp() {

	mType := lorawan.UnconfirmedDataUp
	if d.confirmed {
		mType = lorawan.ConfirmedDataUp
	}

	phy := lorawan.PHYPayload{
		MHDR: lorawan.MHDR{
			MType: mType,
			Major: lorawan.LoRaWANR1,
		},
		MACPayload: &lorawan.MACPayload{
			FHDR: lorawan.FHDR{
				DevAddr: d.devAddr,
				FCnt:    d.fCntUp,
				FCtrl: lorawan.FCtrl{
					ADR: false,
				},
			},
			FPort: &d.fPort,
			FRMPayload: []lorawan.Payload{
				&lorawan.DataPayload{
					Bytes: d.payload,
				},
			},
		},
	}

	if err := phy.EncryptFRMPayload(d.appSKey); err != nil {
		logrus.WithError(err).Error("simulator: encrypt FRMPayload error")
		return
	}

	if err := phy.SetUplinkDataMIC(lorawan.LoRaWAN1_0, 0, 0, 0, d.nwkSKey, d.nwkSKey); err != nil {
		logrus.WithError(err).Error("simulator: set uplink data mic error")
		return
	}
	//d.fCntUp++
	//d.sendUplink(phy)
	//deviceUplinkCounter().Inc()
}