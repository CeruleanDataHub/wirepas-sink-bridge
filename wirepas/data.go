package wirepas

import (
	"fmt"
	"strings"
)

type Message struct {
	DstEP      uint8
	SrcEP      uint8
	SrcAddress uint32
	DstAddress uint32
	Payload    []byte
}

type DestinationEndpoint int

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

	PwsEpDstStatus            DestinationEndpoint = 111
	PwsEpDstTmp               DestinationEndpoint = 112
	PwsEpDstAcc               DestinationEndpoint = 113
	PwsEpDstHumidity          DestinationEndpoint = 114
	PwsEpDstMotion            DestinationEndpoint = 115
	PwsEpDstPressure          DestinationEndpoint = 116
	PwsEpDstIAQ               DestinationEndpoint = 117
	PwsEpDstLight             DestinationEndpoint = 118
	PwsEpDstSound             DestinationEndpoint = 119
	PwsEpDstCurrent           DestinationEndpoint = 120
	PwsEpDstGrideye           DestinationEndpoint = 121
	PwsEpDstGrideyeMultiImage DestinationEndpoint = 122
	PwsEpDstAccBurst          DestinationEndpoint = 123
	PwsEpDstAccShock          DestinationEndpoint = 124
	PwsEpDstButton            DestinationEndpoint = 125
	PwsEpDstADC               DestinationEndpoint = 126
	PwsEpDstAccShockExt       DestinationEndpoint = 127
	PwsEpDstCurrentExt        DestinationEndpoint = 128
	PwsEpDstEnergy            DestinationEndpoint = 129
	PwsEpDstVSAG              DestinationEndpoint = 130
	PwsEpDstCO2               DestinationEndpoint = 131

	PwsEpDstNeighbourRSSI  = 192
	PwsEpDstNeighbourCount = 193
)

func (d *Message) String() string {
	var output strings.Builder
	var s strings.Builder
	for _, b := range d.Payload {
		fmt.Fprintf(&s, "0x%02x ", b)
	}
	fmt.Fprintf(&output, "    DstEp: %d, SrcEp: %d\n    SrcAddr: 0x%02x, DstAddr: 0x%02x\n    Payload: [ %s]\n", d.DstEP, d.SrcEP, d.SrcAddress, d.DstAddress, s.String())

	return output.String()
}
