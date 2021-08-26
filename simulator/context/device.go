package context

import (
	"errors"
	"github.com/brocaar/lorawan"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

type Device struct {
	DevId uint16
	// BanchMark  variables
	Packet_tx, Packet_rx uint8
	PacketLoss int
	Durations DurationSlice
	Start time.Time
	FsmState int
	confirmed bool
	//devEUI lorawan.EUI64
	// AppKey.
	//appKey lorawan.AES128Key
	// Application session-key.
	appSKey lorawan.AES128Key
	// Network session-key.
	nwkSKey lorawan.AES128Key
	// devAddr
	DevAddr lorawan.DevAddr
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
	DownlinkHandleFunc func() error


}

func (d *Device) Marshall() ([]byte, bool){

	phyLoRa, ok := d.UplinkData()
	if !ok {
		log.Fatalf("Error Marshall Phy Lora frame ")
	}
	payload, err := phyLoRa.MarshalBinary()
	if err != nil {
		log.Fatalf("%v", errors.New("Error marshall binary  device Data"))
		return nil, false
	}
	return payload, true
}

func (device *Device) init( id uint16){
	device.DevId = id
	device.DevAddr = counter.getAddr()
}

func (device *Device) GetDevID() (uint16) {
	//b := make([]byte, 2)
	//binary.BigEndian.PutUint16(b, device.DevId)
	//return hex.EncodeToString(b)
	return device.DevId
}

func (device *Device) SetMessagePayload( msg string ){
	device.payload = []byte(msg)
}
func (device *Device) ElapsedTime(){

	if device.FsmState == FSM_RECV {
		t := time.Now()
		elapsed := t.Sub(device.Start)
		device.Durations = append(device.Durations, elapsed)
	} else {
		log.Printf("Not Elapsed time ")
	}
}

// dataUp sends an data uplink.
func (d *Device) UplinkData() (lorawan.PHYPayload, bool) {

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
				DevAddr: d.DevAddr,
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
		return  lorawan.PHYPayload{} , false
	}

	if err := phy.SetUplinkDataMIC(lorawan.LoRaWAN1_0, 0, 0, 0, d.nwkSKey, d.nwkSKey); err != nil {
		logrus.WithError(err).Error("simulator: set uplink data mic error")
		return lorawan.PHYPayload{} , false
	}

	d.fCntUp++

	return phy, true
}

func CreateDevicesForSimulate(devicesLen int){
	if devicesLen < 1 {
		log.Fatalf("Number of devices is not valid")
	}

	for i := 0; i < devicesLen; i++ {
		DevicesContext_Self().NewDevice()
	}
}