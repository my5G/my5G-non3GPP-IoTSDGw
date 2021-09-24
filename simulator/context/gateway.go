package context

import (
	"bytes"
	//"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"sync"
	//"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	//	"log"
	"net"
)

// GatewayOption is the interface for a gateway option.
//type GatewayOption func(*Gateway) error

type Gateway struct {
	Downlink   *net.UDPAddr
	Uplink  *net.UDPAddr
	MAC string

	upLinkMsg *UpStreamJSON
	upMessageMutex sync.Mutex
}


type PushAckOption func(data *PushAck) error

type PushAck struct {
	ProtocolVersion uint8
	RandomToken   uint16
	Identifier uint8
}

func PutPushAckProcotolVersion (payload []byte) PushAckOption {
	return func(p *PushAck) (error){
		if len(payload) < 1 {
			return errors.New("PUSH Ack Protocol Version error, length of byte array less than one byte ")
		}
		value := payload[0]

		//logrus.WithFields( logrus.Fields{
		//	"ProtocolVersion": value,
		//}).Info("Message Push ACK Packet ")


		p.ProtocolVersion = uint8(value)
		return nil
	}
}

func PutPushAckRandomToken (payload []byte) PushAckOption {
	return func(p *PushAck) (error){
		if len(payload) != 4 {
			return errors.New(fmt.Sprintf(" Error protocol Random Token %b \n", payload))
		}

		value := bytes.NewReader( payload[1:3] )
		//u := binary.BigEndian.Uint16(value)

		var u uint16
		binary.Read(value, binary.BigEndian, &u)

		//fmt.Printf(" ************ DEvID = %d \n", u)

		//logrus.WithFields( logrus.Fields{
		//	"ID Token": u,
		//}).Info("Message Push ACK Packet ")

		p.RandomToken = u
		return nil
	}
}

func PutPushAckIdentifier (payload []byte) PushAckOption {
	return func(p *PushAck) (error){
		value := payload[3]

		//logrus.WithFields( logrus.Fields{
		//	"Identifier": value,
		//	}).Info("Message Push ACK Packet ")


		switch value {
		case 1:
			p.Identifier = 1
		default:
			return errors.New(fmt.Sprintf("Error identifier protocol %d \n", value ))
		}
		return nil
	}
}

type PushOption func(data *PushData) error

type PushData struct {
	protocolVersion byte
	randomToken []byte
	Indentifier  byte
	mac []byte
	upData []byte
}

func (p *PushData) Join() ([]byte){
	frame := make([]byte, 0)
	frame = append(p.upData, frame...)
	frame = append(p.mac, frame...)
	frame = append([]byte{p.Indentifier}, frame...)
	frame = append(p.randomToken, frame...)
	frame = append([]byte{p.protocolVersion}, frame...)
	return frame
}

func WithProtocolVersion( version uint8) PushOption {
	return func(p *PushData) error {

		//logrus.WithFields( logrus.Fields{
		//	"Protocol Version": version,
		//}).Info("Message Push Option Packet ")

		/*
		data, err := hex.DecodeString(version)
		if err != nil {
			logrus.WithFields(
				logrus.Fields{
					"Version" : data,
				}).Error("Hexadecimal converter ")
			return errors.Wrap(err, "Error paser protocol version to hexcode byte array ")
		}*/

		p.protocolVersion = version
		return nil
	}
}

func WithRandomToken(token uint16) PushOption {
	return func(p *PushData) error {

		//logrus.WithFields( logrus.Fields{
		//	"Random Token": token,
		//}).Info("Message Push Option Packet ")


		//data, err := hex.DecodeString(token)
		//if err != nil {
		//	logrus.WithFields(
		//		logrus.Fields{
		//			"tokenError" : data,
		//		}).Error("Hexadecimal converter ")
		//	return errors.Wrap(err, "Error parser token to hexcode byte array ")
		//}
		b := make([]byte, 2)
		binary.BigEndian.PutUint16(b, token)
		p.randomToken = b
		return nil
	}
}

func WithIndetifier( id uint8 ) PushOption {
	return func(p *PushData) error {
		//logrus.WithFields( logrus.Fields{
		//	"Indentifier": id,
		//}).Info("Message Push Option Packet ")
		p.Indentifier = id
		return nil
	}
}

func WithMac( mac  string) PushOption {
	return func(p *PushData) error {

		//logrus.WithFields( logrus.Fields{
		//	"Mac address": mac,

		//}).Info("Message Push Option Packet ")

		data, err := hex.DecodeString(mac)
		if err != nil {
			logrus.WithFields(
				logrus.Fields{
					"Mac" : data,
				}).Error("Hexadecimal converter ")
			return errors.Wrap(err, "Error parser MAC to hexcode byte array ")
		}
		p.mac = data
		return nil
	}
}

func WithJsonObject( uplink UpStreamJSON) PushOption {

	return func(g *PushData) error {


		//logrus.WithFields( logrus.Fields{
		//	"Json message": uplink,
		//}).Info("Message Push Option Packet ")


		marshallCode, err := json.Marshal(uplink)
		if err != nil {
			logrus.WithFields(
				logrus.Fields{
					"Json Marshall Message " : uplink,
				}).Error("Error Data convert to json ")
			return errors.Wrap(err, "Error parser JSON UpLink to byte array ")
		}

		g.upData = marshallCode
		return nil
	}
}

func (g *Gateway) NewUplinkEventHandler(loraRFPayload []byte, opts ...PushOption) ([]byte, error) {

	upJson := Rxpk{
		Time: DefaultTime,
		Tmms: DefaultTmms,
		Tmst: DefaultTmst,
		Chan: DefaultChan,
		Rfch: DefaultRfch,
		Freq: DefaultFreq,
		Stat: DefaultStat,
		Modu: DefaultModu,
		Datr: DefaultDatr,
		Codr: DefaultCodr,
		Rssi: DefaultRssi,
		Lsnr: DefaultLsnr,
	}


	data := base64.StdEncoding.EncodeToString(loraRFPayload)
	upJson.Size = uint16( len(loraRFPayload) )
	upJson.Data = data

	uplink := UpStreamJSON{}
	uplink.Rxpk = []Rxpk{upJson}

	pushData := &PushData{}

	for _, o := range opts{
		if err := o(pushData); err != nil {
			return nil, err
		}
	}

	if err := WithJsonObject(uplink)(pushData);  err != nil {
		return nil, err
	}

	//fmt.Printf("%b \n", pushData.Join())
	return pushData.Join(), nil
	//Mount bytes message
}

func (g *Gateway) DownlinkEventHandler(opts ...PushAckOption) (error) {

	pushAck := &PushAck{}

	for _, o := range opts{
		if err := o(pushAck); err != nil  {
			return err
		}
	}

	var id uint16

	/* Load Device*/
	id = pushAck.RandomToken
	device, ok := DevicesContext_Self().DeviceLoad(id)
	if !ok {
		log.Fatalf("Device id not found in recv Dispatch")
		return errors.Errorf("Device id not found in recv Dispatch")
	}

	if device.FsmState != FSM_WAIT {
		return errors.New(fmt.Sprintf("FSM = %d ", FSM_WAIT))
	}

	device.DownlinkHandleFunc = func() error {
		device.Packet_rx++
		device.fCntDown++
		device.ElapsedTime()
		device.FsmState = FSM_RECV
		return nil
	}
	device.DownlinkHandleFunc()
	DevicesContext_Self().Stores.Store(device.DownLinkInfo(true))
	device.DoneRecv <- true
	return nil
}