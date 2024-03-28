package layers

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type IPv4Info struct {
	Source       string `json:"source"`
	Destination  string `json:"destination"`
	PayloadType  string `json:"payload"`
	Length       uint16 `json:"length"`
	HeaderLength int    `json:"header_length"`
}

func (info IPv4Info) Layer() Layer {
	return IPv4Layer
}

func IPv4Constructor(layer gopacket.Layer) Info {
	ipv4Layer := layer.(*layers.IPv4)
	return IPv4Info{
		Source:       ipv4Layer.SrcIP.String(),
		Destination:  ipv4Layer.DstIP.String(),
		Length:       ipv4Layer.Length,
		PayloadType:  ipv4Layer.Protocol.String(),
		HeaderLength: len(ipv4Layer.Contents) - len(ipv4Layer.Payload),
	}
}
