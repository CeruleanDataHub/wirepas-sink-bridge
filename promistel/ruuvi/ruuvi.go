package ruuvi

import (
	"encoding/binary"
	"fmt"

	"github.com/ceruleandatahub/wirepas-sink-bridge/wirepas"
)

type RuuviData struct {
	Temperature float32 `json:"temperature,omitempty"`
	Humidity    float32 `json:"humidity,omitempty"`
	Pressure    uint32  `json:"pressure,omitempty"`
}

func DecodeData(msg *wirepas.Message) (*RuuviData, error) {
	data := &RuuviData{}

	switch wirepas.DestinationEndpoint(msg.DstEP) {
	case wirepas.PwsEpDstTmp:
		data.Temperature = float32(int16(binary.LittleEndian.Uint16(msg.Payload))) / 100
	case wirepas.PwsEpDstHumidity:
		data.Humidity = float32(binary.LittleEndian.Uint16(msg.Payload)) / 100
	case wirepas.PwsEpDstPressure:
		data.Pressure = binary.LittleEndian.Uint32(msg.Payload)
	default:
		return nil, fmt.Errorf("Unsupported measurement type %d", msg.DstEP)
	}

	return data, nil
}
