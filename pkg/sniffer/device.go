package sniffer

import (
	"time"

	"github.com/google/gopacket/pcap"
)

type device string

type DeviceConfig struct {
	name       device
	BufferSize int
	Immediate  bool
	Promisc    bool
	RFMon      bool
	SnapLen    int
	Timeout    time.Duration
}

func NewDeviceConfig(name string) DeviceConfig {
	return DeviceConfig{
		name:       device(name),
		BufferSize: 2097152,
		Immediate:  false,
		Promisc:    false,
		RFMon:      false,
		SnapLen:    262144,
		Timeout:    pcap.BlockForever,
	}
}
