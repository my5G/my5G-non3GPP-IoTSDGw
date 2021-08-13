package udp_server

import (
	"bufio"
	"bytes"
	"github.com/my5G/my5G-non3GPP-IoTSDGw/simulator/context"
	"log"
	"net"
	"sync"
)

var mtx sync.Mutex

var (
	ChannelForward01  chan sendMessage
	ChannelForward02  chan sendMessage
	ChannelForward03  chan sendMessage
	ChannelForward04  chan sendMessage
	ChannelForward05  chan sendMessage
	ChannelForward06  chan sendMessage
	ChannelForward07  chan sendMessage
	ChannelForward08  chan sendMessage
	ChannelForwardRecv chan recvMessage
)

type UDPSendInfoGroup struct  {
	ChannelID int
	Addr *net.UDPAddr
}

type sendMessage struct {
	DstAddr *net.UDPAddr
	Payload []byte
	Length int
}
type recvMessage struct {
	ipSrcAddr string
	Payload []byte
	Length int
}

func init() {
	ChannelForward01 = make(chan sendMessage)
	ChannelForward02 = make(chan sendMessage)
	ChannelForward03 = make(chan sendMessage)
	ChannelForward04 = make(chan sendMessage)
	ChannelForward05 = make(chan sendMessage)
	ChannelForward06 = make(chan sendMessage)
	ChannelForward07 = make(chan sendMessage)
	ChannelForward08 = make(chan sendMessage)
	ChannelForwardRecv = make(chan recvMessage)
}

func RecvMessage(msg recvMessage) {
	mtx.Lock()
	ChannelForwardRecv <- msg
	mtx.Unlock()
}

func Run (){

	link := context.DevicesContext_Self().ForwarderConn

	/*
	listenerPortDownlink, err := net.ListenUDP("udp", link)
	if err != nil {
		log.Fatalf("[IKE] Listen on UDP socket failed: %+v", err)
		return
	}
	*/
	dialUplistenerPort, err := net.DialUDP("udp",nil, link)
	if err != nil {
			log.Fatalf(" Listen on UDP socket failed: %+v", err)
			return
	}

	//Maybe syn all go routines
	//go reader (ChannelIDRecv, listenerPortDownlink)
	go sender (ChannelID1, dialUplistenerPort)
	go sender (ChannelID2, dialUplistenerPort)
	go sender (ChannelID3, dialUplistenerPort)
	go sender (ChannelID4, dialUplistenerPort)
	go sender (ChannelID5, dialUplistenerPort)
	go sender (ChannelID6, dialUplistenerPort)
	go sender (ChannelID7, dialUplistenerPort)
	go sender (ChannelID8, dialUplistenerPort)

}

func reader(chanId int , conn *net.UDPConn){
	reader := bufio.NewReader(conn)
	var buffer bytes.Buffer

	if chanId != ChannelIDRecv {
		log.Fatal("Channel ID out of range Reader")
		return
	}

	for {
			n , err := buffer.ReadFrom(reader)
			if err != nil {
				log.Fatalf("Read from UDP failed: %+v", err)
				return
			}
			if n <= 0 {
				log.Fatalf("Read from UDP failed: %+v", err)
				return
			}

			msg := recvMessage {
			conn.RemoteAddr().String(),
				buffer.Bytes(),
			len(buffer.Bytes()),
		}
		RecvMessage(msg)
	}
}

func sender(channelID int, conn *net.UDPConn) {

	if channelID >= ChannelID1 && channelID <= ChannelID8 {

		for {

			switch channelID {
			case ChannelID1:
				sendData := <-ChannelForward01
				n, err := conn.Write(sendData.Payload)
				if err != nil {
					log.Fatal("Sending data through UDP failed: %+v", err)
				}

				if n != sendData.Length {
					//Make warn
					log.Printf("There is data not being sent\n")
				}

			case ChannelID2:
				sendData := <-ChannelForward02
				n, err := conn.Write(sendData.Payload)
				if err != nil {
					log.Fatal("Sending data through UDP failed: %+v", err)
				}

				if n != sendData.Length {
					//Make warn
					log.Printf("There is data not being sent\n")
				}

			case ChannelID3:
				sendData := <-ChannelForward03
				n, err := conn.Write(sendData.Payload)
				if err != nil {
					log.Fatal("Sending data through UDP failed: %+v", err)
				}

				if n != sendData.Length {
					//Make warn
					log.Printf("There is data not being sent\n")
				}

			case ChannelID4:
				sendData := <-ChannelForward04
				n, err := conn.Write(sendData.Payload)
				if err != nil {
					log.Fatal("Sending data through UDP failed: %+v", err)
				}

				if n != sendData.Length {
					//Make warn
					log.Printf("[IKE] There is data not being sent\n")
				}

			case ChannelID5:
				sendData := <-ChannelForward05
				n, err := conn.Write(sendData.Payload)
				if err != nil {
					log.Fatal("[IKE] Sending data through UDP failed: %+v", err)
				}

				if n != sendData.Length {
					//Make warn
					log.Printf("There is data not being sent\n")
				}

			case ChannelID6:
				sendData := <-ChannelForward06
				n, err := conn.Write(sendData.Payload)
				if err != nil {
					log.Fatal("Sending data through UDP failed: %+v", err)
				}

				if n != sendData.Length {
					//Make warn
					log.Printf("There is data not being sent\n")
				}

			case ChannelID7:
				sendData := <-ChannelForward07
				n, err := conn.Write(sendData.Payload)
				if err != nil {
					log.Fatal("Sending data through UDP failed: %+v", err)
				}

				if n != sendData.Length {
					//Make warn
					log.Printf("There is data not being sent\n")
				}

			case ChannelID8:
				sendData := <-ChannelForward08
				n, err := conn.Write(sendData.Payload)
				if err != nil {
					log.Fatal("Sending data through UDP failed: %+v", err)
				}

				if n != sendData.Length {
					//Make warn
					log.Printf("There is data not being sent\n")
				}
			}
		}
	} else {
		log.Fatal(" Invalid channel ID")
	}
}