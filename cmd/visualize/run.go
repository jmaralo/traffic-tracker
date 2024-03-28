package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"

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

	chart := charts.NewBar()
	chart.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "TCP vs UDP traffic",
		Subtitle: "length with headers",
	}))

	chart.SetXAxis([]string{"TCP", "UDP"}).AddSeries("traffic", generateBars(records))

	output, err := os.Create("renders/tcp_vs_udp.html")
	if err != nil {
		return err
	}
	defer output.Close()

	return chart.Render(output)
}

func generateBars(records []record.Capture) []opts.BarData {
	barData := []opts.BarData{
		{
			Name:  "TCP",
			Value: 0,
		},
		{
			Name:  "UDP",
			Value: 0,
		},
	}

	for _, record := range records {
		if _, ok := record.Layers[layers.TCPLayer]; ok {
			currentLength := barData[0].Value
			barData[0].Value = currentLength.(int) + record.Length
		}
		if _, ok := record.Layers[layers.UDPLayer]; ok {
			currentLength := barData[1].Value
			barData[1].Value = currentLength.(int) + record.Length
		}
	}

	return barData
}
