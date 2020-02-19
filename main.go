package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"main/datasource"
	"main/env"
	"main/sender"
	"main/transport/mcast"
	"main/transport/tcp"

	log "github.com/sirupsen/logrus"
)

func main() {
	config := env.ReadOS()
	config.Dump()

	initLogger(config.LogLevel, config.PrettyLogOutput)

	inputReader, err := datasource.New(config.SourceFile)
	if err != nil {
		fmt.Println("failed to read file:", err.Error())
		return
	}
	inputReader.Dump(os.Stdout)

	wg := &sync.WaitGroup{}

	switch config.OutputMode {
	case "mcast":
		runMcast(wg, inputReader)

	case "tcp_client":
		runTcpClient(wg, inputReader)

	case "tcp_server":
		runTcpServer(wg, inputReader)

	default:
		fmt.Print("unknown mode")
		return
	}

	wg.Wait()
}

func runMcast(wg *sync.WaitGroup, inputReader *datasource.MessageReader) {

	writer := mcast.NewWriter("239.0.112.1:6501")
	sndr := sender.New(writer, inputReader)

	_ = writer.Run()
	sndr.Run()
}

func runTcpClient(wg *sync.WaitGroup, inputReader *datasource.MessageReader) {

	writer := tcp.NewClientWriter("192.168.15.137:8080")
	sndr := sender.New(writer, inputReader)

	_ = writer.Run()
	sndr.Run()
}

func runTcpServer(wg *sync.WaitGroup, inputReader *datasource.MessageReader) {

	onConn := func(w sender.Writer) {
		sndr := sender.New(w, inputReader)
		sndr.Run()
	}

	writer := tcp.NewServerWriter(":8080", onConn)
	_ = writer.Run(wg)
}

func initLogger(logLevel string, pretty bool) {
	if pretty {
		log.SetFormatter(&log.JSONFormatter{})
	}
	log.SetOutput(os.Stderr)

	switch strings.ToLower(logLevel) {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	default:
		log.SetLevel(log.DebugLevel)
	}
}