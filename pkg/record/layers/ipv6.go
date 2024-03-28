package layers

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type IPv6Info struct {
	Source       string `json:"source"`
	Destination  string `json:"destination"`
	PayloadType  string `json:"payload"`
	Length       uint16 `json:"length"`
	HeaderLength int    `json:"header_length"`
}

func (info IPv6Info) Layer() Layer {
	return IPv6Layer
}

func IPv6Constructor(layer gopacket.Layer) Info {
	ipv6Layer := layer.(*layers.IPv6)
	return IPv6Info{
		Source:       ipv6Layer.SrcIP.String(),
		Destination:  ipv6Layer.DstIP.String(),
		Length:       ipv6Layer.Length,
		PayloadType:  ipv6Layer.NextHeader.String(),
		HeaderLength: len(ipv6Layer.Contents) - len(ipv6Layer.Payload),
	}
}
