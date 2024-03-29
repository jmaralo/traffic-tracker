package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/google/gopacket/pcap"
	"github.com/jmaralo/traffic-tracker/pkg/sniffer"
)

var promiscFlag = flag.Bool("p", false, "run in promisc mode")

func main() {
	flag.Parse()

	devices, err := pcap.FindAllDevs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	scanner := bufio.NewReader(os.Stdin)
	configs := make([]sniffer.DeviceConfig, 0, len(devices))
	for _, device := range devices {
		include, err := promptYesNo(scanner, fmt.Sprintf("Include device %s?", device.Name))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if !include {
			continue
		}

		config := sniffer.NewDeviceConfig(device.Name)
		config.Promisc = *promiscFlag
		configs = append(configs, config)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	err = Run(configs, signalChan)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}

func promptYesNo(scanner *bufio.Reader, prompt string) (bool, error) {
	fmt.Printf("%s [Y/N] > ", prompt)
	for {
		answer, err := scanner.ReadString('\n')
		if err != nil {
			return false, err
		}
		answer = strings.ToLower(answer)
		if answer[0] == 'y' {
			return true, nil
		} else if answer[0] == 'n' {
			return false, nil
		}

		fmt.Print("please input [Y/N] > ")
	}
}
