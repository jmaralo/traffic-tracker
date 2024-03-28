package layers

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type EthernetInfo struct {
	Source       string `json:"source"`
	Destination  string `json:"destination"`
	PayloadType  string `json:"payload"`
	Length       uint16 `json:"length"`
	HeaderLength int
}

func (info EthernetInfo) Layer() Layer {
	return EthernetLayer
}

func EthernetConstructor(layer gopacket.Layer) Info {
	ethernetLayer := layer.(*layers.Ethernet)
	return EthernetInfo{
		Source:       ethernetLayer.SrcMAC.String(),
		Destination:  ethernetLayer.DstMAC.String(),
		PayloadType:  ethernetLayer.EthernetType.String(),
		Length:       ethernetLayer.Length,
		HeaderLength: len(ethernetLayer.Contents) - len(ethernetLayer.Payload),
	}
}
