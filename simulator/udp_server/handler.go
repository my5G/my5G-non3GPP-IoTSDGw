package udp_server

import (
	"github.com/my5G/my5G-non3GPP-IoTSDGw/simulator/context"
	"github.com/sirupsen/logrus"
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
	Value       interface{}
}

func CycleChannel() int {
	reqChannelNumberMutex.Lock()
		channelFlag++
		if channelFlag > 8 { channelFlag = 1 }
	reqChannelNumberMutex.Unlock()
	return channelFlag
}

func SendChannelMessage(packet []byte, tokenId uint16, channelID int) {

	gateway := context.DevicesContext_Self().Gateway

	message, err := gateway.NewUplinkEventHandler(
		packet,
		context.WithProtocolVersion(2),
		context.WithRandomToken(tokenId),
		context.WithIndetifier(0),
		context.WithMac(gateway.MAC),
		)

	if err != nil {
		log.Fatalf("Packet gateway channel %d make errror %v", channelID, err)
	}

	chanMsg := sendMessage{
		message,
		len(message),
	}

	switch channelID {
	case 1:
		ChannelForward01 <- chanMsg
	case 2:
		ChannelForward02 <- chanMsg
	case 3:
		ChannelForward03 <- chanMsg
	case 4:
		ChannelForward04 <- chanMsg
	case 5:
		ChannelForward05 <- chanMsg
	case 6:
		ChannelForward06 <- chanMsg
	case 7:
		ChannelForward07 <- chanMsg
	case 8:
		ChannelForward08 <- chanMsg
	default:
		log.Printf("Channel ID #{channelID} not found")
	}
}

func HandleRecvMessage(){
	for {
		select {
		case msg, ok := <-ChannelForwardRecv:
			if ok {
				go Dispatch(msg.Payload)
				//a, _ := context.DevicesContext_Self().DeviceLoad(0)
			}
		}
	}
}

func Dispatch(payload []byte){

	if len(payload) != 4 {
		log.Fatalf("ACK Recv payload Error")
	}

	err := context.DevicesContext_Self().Gateway.DownlinkEventHandler(
		context.PutPushAckProcotolVersion(payload),
		context.PutPushAckRandomToken(payload),
		context.PutPushAckIdentifier(payload),
		)

	if err != nil {
		logrus.Errorf("Error recv %v", err )
	}

}