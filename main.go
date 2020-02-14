package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"

	"main/datasource"
	"main/transport"
	"main/transport/mcast"
	"main/transport/tctip"
)

func main() {

	inputReader, err := datasource.New("data-multi-s4.json")
	if err != nil {
		fmt.Println("failed to read file:", err.Error())
		return
	}

	inputReader.Dump(os.Stdout)

	writer := mcast.NewWriter("239.0.112.1:6501")
	sender := transport.NewSender(writer, inputReader)

	writer.Run()
	sender.Run()

}

func RunTcpClient(ctx context.Context, wg *sync.WaitGroup, inReader *datasource.MessageReader, cancel context.CancelFunc) {
	//addr := "192.168.15.137:8080"
	addr := "127.0.0.1:8080"

	tctip.StartClient(ctx, wg, addr, "0", inReader, cancel)
}

func RunTcpServer(ctx context.Context, wg *sync.WaitGroup, inReader *datasource.MessageReader, cancel context.CancelFunc) {
	//addr := "192.168.15.137:8080"
	addr := "127.0.0.1:8080"

	tctip.StartServer(ctx, wg, addr, inReader, cancel)
}



func addresses() {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("failed to get interfaces")
		return
	}

	for _, i := range ifaces {
		fmt.Println()
		fmt.Println(i.Name)

		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println("failed to get addresses")
		} else {
			fmt.Println("\tlocal addresses")
			for _, a := range addrs {
				fmt.Printf("\t: %s\n", a.String())
			}
		}

		fmt.Println()

		mcaddrs, err := i.MulticastAddrs()
		if err != nil {
			fmt.Println("failed to get addresses")
		} else {
			fmt.Println("\tmulticast addresses")
			for _, a := range mcaddrs {
				fmt.Printf("\t: %s\n", a.String())
			}
		}

		fmt.Println()

	}
}

