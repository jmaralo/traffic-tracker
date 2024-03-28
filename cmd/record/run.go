package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/jmaralo/traffic-tracker/pkg/sniffer"
)

func Run(configs []sniffer.DeviceConfig, signalChan <-chan os.Signal) error {
	sniffer, err := sniffer.New(configs, 100)
	if err != nil {
		return err
	}
	defer sniffer.Close()

	saveFile, err := os.Create(path.Join("records", fmt.Sprintf("traffic_%d.json", time.Now().Unix())))
	if err != nil {
		return err
	}
	defer saveFile.Close()

	written := int64(0)
	defer fmt.Printf("Written a total of %d bytes\n", written)
	reminderTicker := time.NewTicker(time.Minute)
	for {
		select {
		case record, ok := <-sniffer.Records():
			if !ok {
				return nil
			}
			write, err := saveJSON(record, saveFile)
			if err != nil {
				return err
			}
			written += write
		case <-signalChan:
			return nil
		case <-reminderTicker.C:
			fmt.Printf("Written a total of %d bytes\n", written)
		}
	}
}

func saveJSON(data any, file *os.File) (int64, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	return io.Copy(file, bytes.NewReader(raw))
}
