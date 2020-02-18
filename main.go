package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"main/datasource"
	"main/transport/mcast"
	"main/transport/sender"
	"main/transport/tcp"
)

func main() {

	mode := flag.String("mode", "unknown", "input mode")
	flag.Parse()

	inputReader, err := datasource.New("data-multi-s2.json")
	if err != nil {
		fmt.Println("failed to read file:", err.Error())
		return
	}
	inputReader.Dump(os.Stdout)

	wg := &sync.WaitGroup{}

	switch *mode {
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
	writer := tcp.NewServerWriter(":8080", inputReader)
	_ = writer.Run(wg)
}
