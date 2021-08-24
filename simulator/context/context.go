package context

import (
	"fmt"
	"log"
	"net"
	"sync"
)

var ctx DevContext
//var Gw GwMessage

type DevContext struct{
	DevicesPool  sync.Map
	Gateway *Gateway
}

func DevicesContext_Self() *DevContext{
	return &ctx
}

func ( ctx *DevContext) NewDevice() *Device {

	valueID := Incremment()

	if  valueID < 0  {
		log.Fatalf("Dev Id code Error")
	}

	device := new(Device)
	device.init(valueID)
	device.Durations = make(DurationSlice, 0, 1000)
	device.nwkSKey = NWKSKEYTestOnly
	device.appSKey = APPSKEYTestOnly
	device.fPort = 2 // Make
	//device.payload = []byte(Message)
	device.FsmState = FSM_IDLE

	ctx.DevicesPool.Store(valueID, device)

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

func ( ctx *DevContext) ConfigSocketUDPAddr( ipAddr string, port int)(bool){
	serverAddr,err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d",
		ipAddr, port))
	if err != nil {
		log.Fatalf("Config Bind IPAddr Error %v", err)
		return false
	}
	ctx.Gateway.Downlink = serverAddr

	return true
}

func ( ctx *DevContext) ConfigUplink( ipAddr string, port int)(bool){
	serverAddr,err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d",
		ipAddr, port))
	if err != nil {
		log.Fatalf("Config Bind IPAddr Error %v", err)
		return false
	}
	ctx.Gateway.Uplink = serverAddr

	return true
}