package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var logPath = flag.String("i", "", "Path to log file")

func main() {
	flag.Parse()

	if *logPath == "" {
		fmt.Fprintln(os.Stderr, "Missing log file (-i)")
		os.Exit(1)
	}

	file, err := os.Open(*logPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	err = Run(file, signalChan)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}
