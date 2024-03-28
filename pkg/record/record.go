package record

import (
	"encoding/json"
	"time"

	"github.com/google/gopacket"
	_ "github.com/google/gopacket/layers"
	"github.com/jmaralo/traffic-tracker/pkg/record/layers"
)

type Capture struct {
	Device    string                       `json:"device"`
	Length    int                          `json:"length"`
	Timestamp time.Time                    `json:"timestamp"`
	Layers    map[layers.Layer]layers.Info `json:"layers"`
}

func New(packet gopacket.Packet) Capture {
	metadata := packet.Metadata()
	return Capture{
		Length:    metadata.Length,
		Timestamp: metadata.Timestamp,
		Layers:    layers.Interpret(packet.Layers()),
	}
}

type partialCapture struct {
	Device    string                           `json:"device"`
	Length    int                              `json:"length"`
	Timestamp time.Time                        `json:"timestamp"`
	Layers    map[layers.Layer]json.RawMessage `json:"layers"`
}

func (capture *Capture) UnmarshalJSON(data []byte) error {
	var partial partialCapture
	err := json.Unmarshal(data, &partial)
	if err != nil {
		return err
	}

	capture.Device = partial.Device
	capture.Length = partial.Length
	capture.Timestamp = partial.Timestamp

	completeLayers := make(map[layers.Layer]layers.Info, len(partial.Layers))
	for layer, layerData := range partial.Layers {
		var err error = nil
		switch layer {
		case layers.EthernetLayer:
			var ethernet layers.EthernetInfo
			err = json.Unmarshal(layerData, &ethernet)
			completeLayers[layer] = ethernet
		case layers.IPv4Layer:
			var ipv4 layers.IPv4Info
			err = json.Unmarshal(layerData, &ipv4)
			completeLayers[layer] = ipv4
		case layers.IPv6Layer:
			var ipv6 layers.IPv6Info
			err = json.Unmarshal(layerData, &ipv6)
			completeLayers[layer] = ipv6
		case layers.TCPLayer:
			var tcp layers.TCPInfo
			err = json.Unmarshal(layerData, &tcp)
			completeLayers[layer] = tcp
		case layers.UDPLayer:
			var udp layers.UDPInfo
			err = json.Unmarshal(layerData, &udp)
			completeLayers[layer] = udp
		}
		if err != nil {
			return err
		}
	}
	capture.Layers = completeLayers

	return nil
}
