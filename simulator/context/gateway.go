package context

import (
	//"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/pkg/errors"
	"sync"

	//"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	//	"log"
	"net"
)

// GatewayOption is the interface for a gateway option.
//type GatewayOption func(*Gateway) error

type PushOption func(data *PushData) error

type PushData struct {
	protocolVersion []byte
	randomToken []byte
	Indentifier  uint8
	mac []byte
	upData []byte
}

type Gateway struct {
	Downlink   *net.UDPAddr
	Uplink  *net.UDPAddr
	MAC string

	upLinkMsg *UpStreamJSON
	upMessageMutex sync.Mutex
}

func (p *PushData) Join() ([]byte){
	frame := make([]byte, 0)
	frame = append(p.protocolVersion, frame...)
	frame = append(p.randomToken, frame...)
	frame = append([]byte{p.Indentifier}, frame...)
	frame = append(p.mac, frame...)
	frame = append(p.upData, frame...)
	return frame
}

func WithProtocolVersion( version string) PushOption{
	return func(p *PushData) error {

		logrus.WithFields( logrus.Fields{
			"Protocol Version": version,
		}).Info("Message Push Option Packet ")

		data, err := hex.DecodeString(version)
		if err != nil {
			logrus.WithFields(
				logrus.Fields{
					"Version" : data,
				}).Error("Hexadecimal converter ")
			return errors.Wrap(err, "Error paser protocol version to hexcode byte array ")
		}

		p.protocolVersion = data
		return nil
	}
}

func WithRandomToken(token string) PushOption{
	return func(p *PushData) error {

		logrus.WithFields( logrus.Fields{
			"Random Token": token,
		}).Info("Message Push Option Packet ")

		data, err := hex.DecodeString(token)
		if err != nil {
			logrus.WithFields(
				logrus.Fields{
					"tokenError" : data,
				}).Error("Hexadecimal converter ")
			return errors.Wrap(err, "Error parser token to hexcode byte array ")
		}

		p.randomToken = data
		return nil
	}
}

func WithIndetifier( id uint8 ) PushOption{
	return func(p *PushData) error {
		logrus.WithFields( logrus.Fields{
			"Indentifier": id,
		}).Info("Message Push Option Packet ")
		p.Indentifier = id
		return nil
	}
}

func WithMac( mac string) PushOption{
	return func(p *PushData) error {

		logrus.WithFields( logrus.Fields{
			"Mac address": mac,

		}).Info("Message Push Option Packet ")

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

func WithJsonObject( uplink UpStreamJSON) PushOption{

	return func(g *PushData) error {

		logrus.WithFields( logrus.Fields{
			"Json message": uplink,
		}).Info("Message Push Option Packet ")

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

	return pushData.Join(), nil
	//Mount bytes message
}

//func (g *Gateway) downlinkEventHandler(c mqtt.Client, msg mqtt.Message) {
//}