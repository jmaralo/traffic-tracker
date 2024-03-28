package sniffer

import (
	"fmt"

	"github.com/google/gopacket"
	_ "github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/jmaralo/traffic-tracker/pkg/record"
)

type handler struct {
	deviceName string
	capture    *pcap.Handle
	source     *gopacket.PacketSource
	recordChan chan record.Capture
}

func handle(device string, capture *pcap.Handle, recordBuffer int) handler {
	fmt.Println(capture.LinkType())
	handler := handler{
		deviceName: device,
		capture:    capture,
		source:     gopacket.NewPacketSource(capture, capture.LinkType()),
		recordChan: make(chan record.Capture, recordBuffer),
	}

	go handler.listen()

	return handler
}

func (handler *handler) listen() {
	defer close(handler.recordChan)

	for packet := range handler.source.Packets() {
		packetRecord := record.New(packet)
		packetRecord.Device = handler.deviceName
		handler.recordChan <- packetRecord
	}
}

func (handler *handler) Records() <-chan record.Capture {
	return handler.recordChan
}

func (handler *handler) Close() {
	handler.capture.Close()
}
