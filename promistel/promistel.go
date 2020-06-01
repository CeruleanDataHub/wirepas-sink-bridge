package promistel

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/houston-inc/wirepas-sink-bridge/wirepas"
)

type MeasurementType int

const (
	// Source endpoints
	PwsEpSrcCommand   = 10
	PwsEpSrcPws       = 110
	PwsEpSrcPromistel = 111
	PwsEpSrcRuuvi     = 112
	PwsEpSrcHaltian   = 113
	PwsEpSrcNordic    = 114
	PwsEpSrcUBlox     = 115

	// Destination endpoints
	PwsEpDstDebug     = 1
	PwsEpDstError     = 2
	PwsEpDstPackDelay = 3
	PwsEpDstPackTime  = 4

	PwsEpDstCommand = 10
	PwsEpDstIntSet  = 13

	PwsEpDstStatus            MeasurementType = 111
	PwsEpDstTmp               MeasurementType = 112
	PwsEpDstAcc               MeasurementType = 113
	PwsEpDstHumidity          MeasurementType = 114
	PwsEpDstMotion            MeasurementType = 115
	PwsEpDstPressure          MeasurementType = 116
	PwsEpDstIAQ               MeasurementType = 117
	PwsEpDstLight             MeasurementType = 118
	PwsEpDstSound             MeasurementType = 119
	PwsEpDstCurrent           MeasurementType = 120
	PwsEpDstGrideye           MeasurementType = 121
	PwsEpDstGrideyeMultiImage MeasurementType = 122
	PwsEpDstAccBurst          MeasurementType = 123
	PwsEpDstAccShock          MeasurementType = 124
	PwsEpDstButton            MeasurementType = 125
	PwsEpDstADC               MeasurementType = 126
	PwsEpDstAccShockExt       MeasurementType = 127
	PwsEpDstCurrentExt        MeasurementType = 128
	PwsEpDstEnergy            MeasurementType = 129
	PwsEpDstVSAG              MeasurementType = 130
	PwsEpDstCO2               MeasurementType = 131

	PwsEpDstNeighbourRSSI  = 192
	PwsEpDstNeighbourCount = 193
)

type Device struct {
	Address string `json:"address"`
}

type RuuviData struct {
	Temperature float32 `json:"temperature,omitempty"`
	Humidity    float32 `json:"humidity,omitempty"`
	Pressure    uint32  `json:"pressure,omitempty"`
}

type DeviceInfo struct {
	Device Device     `json:"device"`
	Rssi   int8       `json:"rssi"`
	Values *RuuviData `json:"sensors"`
}

func (i *DeviceInfo) JSON() (string, error) {
	bs, err := json.Marshal(i)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

func DecodeWirepasMessage(msg *wirepas.Message) (*DeviceInfo, error) {
	device := Device{
		Address: fmt.Sprintf("0x%02x", msg.SrcAddress),
	}

	v, err := DecodeRuuviTagMessage(msg)
	if err != nil {
		return nil, err
	}

	info := &DeviceInfo{
		Device: device,
		Rssi:   0,
		Values: v,
	}
	return info, nil
}

func DecodeRuuviTagMessage(msg *wirepas.Message) (*RuuviData, error) {
	data := &RuuviData{}

	switch MeasurementType(msg.DstEP) {
	case PwsEpDstTmp:
		data.Temperature = float32(int16(binary.LittleEndian.Uint16(msg.Payload))) / 100
	case PwsEpDstHumidity:
		data.Humidity = float32(binary.LittleEndian.Uint16(msg.Payload)) / 100
	case PwsEpDstPressure:
		data.Pressure = binary.LittleEndian.Uint32(msg.Payload)
	default:
		return nil, fmt.Errorf("Unsupported measurement type %d", msg.DstEP)
	}

	return data, nil
}
