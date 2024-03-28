package layers

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type Layer string

const (
	EthernetLayer Layer = "eth"
	IPv4Layer     Layer = "ipv4"
	IPv6Layer     Layer = "ipv6"
	TCPLayer      Layer = "tcp"
	UDPLayer      Layer = "udp"
)

var gopacketToSelf = map[gopacket.LayerType]Layer{
	layers.LayerTypeEthernet: EthernetLayer,
	layers.LayerTypeIPv4:     IPv4Layer,
	layers.LayerTypeIPv6:     IPv6Layer,
	layers.LayerTypeTCP:      TCPLayer,
	layers.LayerTypeUDP:      UDPLayer,
}

type Info interface {
	Layer() Layer
}

type infoConstructor = func(gopacket.Layer) Info

var LayerHandlers map[Layer]infoConstructor = map[Layer]infoConstructor{
	EthernetLayer: EthernetConstructor,
	IPv4Layer:     IPv4Constructor,
	IPv6Layer:     IPv6Constructor,
	TCPLayer:      TCPConstructor,
	UDPLayer:      UDPConstructor,
}

func Interpret(layers []gopacket.Layer) map[Layer]Info {
	layerInfo := make(map[Layer]Info, len(layers))
	for _, layer := range layers {
		selfLayer, ok := gopacketToSelf[layer.LayerType()]
		if !ok {
			continue
		}
		handler, ok := LayerHandlers[selfLayer]
		if !ok {
			continue
		}

		layerInfo[selfLayer] = handler(layer)
	}
	return layerInfo
}
