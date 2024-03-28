package layers

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type UDPInfo struct {
	Source       uint16 `json:"source"`
	Destination  uint16 `json:"destination"`
	HeaderLength int    `json:"header_length"`
}

func (info UDPInfo) Layer() Layer {
	return UDPLayer
}

func UDPConstructor(layer gopacket.Layer) Info {
	tcpLayer := layer.(*layers.UDP)
	return UDPInfo{
		Source:       uint16(tcpLayer.SrcPort),
		Destination:  uint16(tcpLayer.DstPort),
		HeaderLength: len(tcpLayer.Contents) - len(tcpLayer.Payload),
	}
}
