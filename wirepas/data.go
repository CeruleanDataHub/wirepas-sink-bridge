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

func (d *Message) String() string {
	var output strings.Builder
	var s strings.Builder
	for _, b := range d.Payload {
		fmt.Fprintf(&s, "0x%02x ", b)
	}
	fmt.Fprintf(&output, "    DstEp: %d, SrcEp: %d\n    SrcAddr: 0x%02x, DstAddr: 0x%02x\n    Payload: [ %s]\n", d.DstEP, d.SrcEP, d.SrcAddress, d.DstAddress, s.String())

	return output.String()
}
