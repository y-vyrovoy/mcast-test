package main

import (
	"fmt"
	"net"
	"os"
	"sync"

	"main/datasource"
	"main/transport/tcp"
)

func main() {

	inputReader, err := datasource.New("data-multi-s2.json")
	if err != nil {
		fmt.Println("failed to read file:", err.Error())
		return
	}

	inputReader.Dump(os.Stdout)

	wg := &sync.WaitGroup{}

	// Multicast
	//sender := mcast.NewWriter("239.0.112.1:6501")

	// tcp client
	//sender := tcp.NewClientWriter("192.168.15.137:8080")

	// tcp server
	writer := tcp.NewServerWriter(":8080", inputReader)

	//sender := transport.NewSender(sender, inputReader)

	writer.Run(wg)
	//sender.Run()


	wg.Wait()
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

