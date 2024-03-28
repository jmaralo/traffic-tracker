package sniffer

import (
	"fmt"

	"github.com/google/gopacket/pcap"
	"github.com/jmaralo/traffic-tracker/pkg/record"
)

type Sniffer struct {
	handlers   map[device]handler
	recordChan chan record.Capture
}

func New(devices []DeviceConfig, recordBuffer int) (Sniffer, error) {
	sniffer := Sniffer{
		handlers:   make(map[device]handler, len(devices)),
		recordChan: make(chan record.Capture, recordBuffer),
	}

	handlerRecordBuffer := recordBuffer / len(devices)
	for _, config := range devices {
		capture, err := createNewCapture(config)
		if err != nil {
			return sniffer, err
		}

		handler := handle(string(config.name), capture, handlerRecordBuffer)
		sniffer.handlers[config.name] = handler
		go sniffer.addRecordSource(handler.Records())
	}

	return sniffer, nil
}

func createNewCapture(config DeviceConfig) (*pcap.Handle, error) {
	inactiveHandle, err := pcap.NewInactiveHandle(string(config.name))
	defer inactiveHandle.CleanUp()
	if err != nil {
		return nil, err
	}

	err = inactiveHandle.SetBufferSize(config.BufferSize)
	if err != nil {
		fmt.Printf("failed to set buffer size for %s\n", config.name)
	}
	err = inactiveHandle.SetImmediateMode(config.Immediate)
	if err != nil {
		fmt.Printf("failed to set immediate mode for %s\n", config.name)
	}
	err = inactiveHandle.SetPromisc(config.Promisc)
	if err != nil {
		fmt.Printf("failed to set promisc mode for %s\n", config.name)
	}
	err = inactiveHandle.SetRFMon(config.RFMon)
	if err != nil {
		fmt.Printf("failed to set rf monitor mode for %s\n", config.name)
	}
	err = inactiveHandle.SetSnapLen(config.SnapLen)
	if err != nil {
		fmt.Printf("failed to set snaplen for %s\n", config.name)
	}
	err = inactiveHandle.SetTimeout(config.Timeout)
	if err != nil {
		fmt.Printf("failed to set timeout for %s\n", config.name)
	}

	return inactiveHandle.Activate()
}

func (sniffer *Sniffer) addRecordSource(records <-chan record.Capture) {
	for record := range records {
		sniffer.recordChan <- record
	}
}

func (sniffer *Sniffer) Records() <-chan record.Capture {
	return sniffer.recordChan
}

func (sniffer *Sniffer) Close() {
	for _, handler := range sniffer.handlers {
		handler.Close()
	}
}
