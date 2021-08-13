package udp_server

import (
	"encoding/json"
	"fmt"
	"github.com/my5G/my5G-non3GPP-IoTSDGw/simulator/context"
	"log"
	"sync"
)

var (
	channelFlag int
	reqChannelNumberMutex sync.Mutex
	sendChannelMutex sync.Mutex
)

type HandlerMessage struct {
	Event       Event
	UDPSendInfo *UDPSendInfoGroup // used only when Event == EventN1UDPMessage
	Value  interface{}
}

func CycleChannel() int {
	reqChannelNumberMutex.Lock()
		channelFlag++
		if channelFlag > 8 { channelFlag = 1 }
	reqChannelNumberMutex.Unlock()
	return channelFlag
}

func SendChannelMessage(msg []byte, channelID int) {

	chanMsg := sendMessage {
		context.DevicesContext_Self().ForwarderConn,
		msg,
		len(msg),
	}

	switch channelID {
	case 1:
		ChannelForward01  <- chanMsg
	case 2:
		ChannelForward02  <- chanMsg
	case 3:
		ChannelForward03  <- chanMsg
	case 4:
		ChannelForward04  <- chanMsg
	case 5:
		ChannelForward05  <- chanMsg
	case 6:
		ChannelForward06  <- chanMsg
	case 7:
		ChannelForward07  <- chanMsg
	case 8:
		ChannelForward08  <- chanMsg
	default:
		log.Printf("Channel ID #{channelID} not found")
	}
}

func HandleRecvMessage(){
	for {
		select {
		case msg, ok := <- ChannelForwardRecv:
			if ok {
				fmt.Printf("Received packet %s\n", msg.Payload )
				go Dispatch(msg.Payload)
				//a, _ := context.DevicesContext_Self().DeviceLoad(0)
			}
		}
	}
}

func Dispatch( packet []byte){

	var code string
	if err := json.Unmarshal(packet, &code); err != nil {
		log.Fatal("Error Decode Json Packet")
		return
	}

	fmt.Printf("%s\n", code )
	// Vem como 04 hexadecimal
	var id uint16

	/* Load Device*/
	device, ok := context.DevicesContext_Self().DeviceLoad(id)
	if !ok {
		log.Panicf("Device id not found in recv Dispatch")
		return
	}

	device.FsmState = context.FSM_RECV
	device.Packet_rx++

}
