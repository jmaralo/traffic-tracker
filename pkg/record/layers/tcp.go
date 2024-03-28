package layers

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type TCPInfo struct {
	Source       uint16 `json:"source"`
	Destination  uint16 `json:"destination"`
	HeaderLength int    `json:"header_length"`
}

func (info TCPInfo) Layer() Layer {
	return TCPLayer
}

func TCPConstructor(layer gopacket.Layer) Info {
	tcpLayer := layer.(*layers.TCP)
	return TCPInfo{
		Source:       uint16(tcpLayer.SrcPort),
		Destination:  uint16(tcpLayer.DstPort),
		HeaderLength: len(tcpLayer.Contents) - len(tcpLayer.Payload),
	}
}
