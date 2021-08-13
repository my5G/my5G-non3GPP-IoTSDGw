package context

import (
	"encoding/hex"
	"fmt"
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
