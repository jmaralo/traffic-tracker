package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"slices"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/jmaralo/traffic-tracker/pkg/record"
	"github.com/jmaralo/traffic-tracker/pkg/record/layers"
)

func Run(input io.Reader, signalChan <-chan os.Signal) error {
	decoder := json.NewDecoder(input)

	records := make([]record.Capture, 0)
	var nextRecord record.Capture
	for decoder.More() {
		err := decoder.Decode(&nextRecord)
		if err == nil {
			records = append(records, nextRecord)
			continue
		}

		if errors.Is(err, io.EOF) {
			break
		}

		return err
	}

	hosts := getHosts(records)

	chart := charts.NewBar()
	chart.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Traffic per host",
		Subtitle: "bytes, TCP vs UDP",
	}), charts.WithDataZoomOpts(opts.DataZoom{
		Type: "inside",
	}))

	chart.SetXAxis(getHostsLabels(hosts)).AddSeries("TCP", generateBarsTCP(hosts, records)).AddSeries("UDP", generateBarsUDP(hosts, records))

	output, err := os.Create("renders/host_traffic.html")
	if err != nil {
		return err
	}
	defer output.Close()

	fmt.Println("Done!")
	fmt.Println("Data summary:")
	fmt.Printf("\tNumber of different hosts: %d\n", len(hosts))
	fmt.Printf("\tTotal data received: %d bytes\n", getTotalData(records))
	fmt.Printf("\tScan duration %v\n", getScanDuration(records))
	return chart.Render(output)
}

func getHosts(records []record.Capture) map[string]string {
	fmt.Println("getting hosts")
	hosts := make(map[string]string)
	for _, record := range records {
		layer, ok := record.Layers[layers.IPv4Layer]
		if !ok {
			continue
		}
		ipv4 := layer.(layers.IPv4Info)

		if _, hit := hosts[ipv4.Source]; hit {
			continue
		}

		fmt.Printf("\tgetting %s\n", ipv4.Source)
		names, err := net.LookupAddr(ipv4.Source)
		if err != nil {
			hosts[ipv4.Source] = ipv4.Source
		} else {
			hosts[ipv4.Source] = names[0]
		}
	}

	return hosts
}

func getHostsLabels(hosts map[string]string) []string {
	labels := make([]string, 0, len(hosts))
	for _, name := range hosts {
		labels = append(labels, name)
	}
	return labels
}

func generateBarsTCP(hosts map[string]string, records []record.Capture) []opts.BarData {
	fmt.Println("Generating TCP bars")
	hostOrder := getHostsOrder(hosts)
	orderMap := make(map[string]int, len(hosts))
	bars := make([]opts.BarData, len(hosts))
	for i := 0; i < len(hosts); i++ {
		bars[i] = opts.BarData{
			Name:  hosts[hostOrder[i]],
			Value: 0,
			Tooltip: &opts.Tooltip{
				Show:      true,
				Trigger:   "item",
				TriggerOn: "mousemove|click",
			},
		}
		orderMap[hostOrder[i]] = i
	}

	for _, record := range records {
		layer, ok := record.Layers[layers.IPv4Layer]
		if !ok {
			continue
		}
		ipv4 := layer.(layers.IPv4Info)

		if _, ok := record.Layers[layers.TCPLayer]; !ok {
			continue
		}

		i, ok := orderMap[ipv4.Source]
		if !ok {
			continue
		}

		currentValue := bars[i].Value
		bars[i].Value = record.Length + currentValue.(int)
	}

	return bars
}

func generateBarsUDP(hosts map[string]string, records []record.Capture) []opts.BarData {
	fmt.Println("Generating UDP bars")
	hostOrder := getHostsOrder(hosts)
	orderMap := make(map[string]int, len(hosts))
	bars := make([]opts.BarData, len(hosts))
	for i := 0; i < len(hosts); i++ {
		bars[i] = opts.BarData{
			Name:  hosts[hostOrder[i]],
			Value: 0,
			Tooltip: &opts.Tooltip{
				Show:      true,
				Trigger:   "item",
				TriggerOn: "mousemove|click",
			},
		}
		orderMap[hostOrder[i]] = i
	}

	for _, record := range records {
		layer, ok := record.Layers[layers.IPv4Layer]
		if !ok {
			continue
		}
		ipv4 := layer.(layers.IPv4Info)

		if _, ok := record.Layers[layers.UDPLayer]; !ok {
			continue
		}

		i, ok := orderMap[ipv4.Source]
		if !ok {
			continue
		}

		currentValue := bars[i].Value
		bars[i].Value = record.Length + currentValue.(int)
	}

	return bars
}

func getHostsOrder(hosts map[string]string) []string {
	hostList := make([]string, 0, len(hosts))
	for host := range hosts {
		hostList = append(hostList, host)
	}
	slices.Sort(hostList)

	return hostList
}

func getTotalData(records []record.Capture) int {
	total := 0
	for _, record := range records {
		total += record.Length
	}
	return total
}

func getScanDuration(records []record.Capture) time.Duration {
	if len(records) == 0 {
		return 0
	}

	var first time.Time = records[0].Timestamp
	var last time.Time = records[0].Timestamp
	for _, record := range records {
		if first.After(record.Timestamp) {
			first = record.Timestamp
		}
		if last.Before(record.Timestamp) {
			last = record.Timestamp
		}
	}

	return last.Sub(first)
}
