package promistel

import (
	"encoding/json"
	"fmt"

	"github.com/houston-inc/wirepas-sink-bridge/promistel/ruuvi"
	"github.com/houston-inc/wirepas-sink-bridge/wirepas"
)

type Device struct {
	Address string `json:"address"`
}

type DeviceInfo struct {
	Device Device      `json:"device"`
	Rssi   int8        `json:"rssi"`
	Values interface{} `json:"sensors"`
}

func (i *DeviceInfo) JSON() (string, error) {
	bs, err := json.Marshal(i)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

func DecodeMessage(msg *wirepas.Message) (*DeviceInfo, error) {
	device := Device{
		Address: fmt.Sprintf("0x%02x", msg.SrcAddress),
	}

	v, err := ruuvi.DecodeData(msg)
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
