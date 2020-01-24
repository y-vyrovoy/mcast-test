package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"sync"

	"github.com/pkg/errors"
	"main/tctip"
)

func main() {

	//addresses()

	wg := &sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())

	//addr := "224.0.0.1:8888"
	//addr := "239.0.112.1:6501"

	//netInterface := "en0"

	//mcast.StartServer(ctx, wg, netInterface, addr, "0", cancel)
	//mcast.StartServer(ctx, wg, netInterface, addr, "1", cancel)
	//mcast.StartPing(ctx, wg, addr, "3", cancel)

	tctip.StartServer(ctx, wg, ":8080", "0", cancel)
	//tctip.StartServer(ctx, wg, ":8080", "1", cancel)

	tctip.StartPing(ctx, wg, ":8080", "0", cancel)



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

func startSendNmea(ctx context.Context, wg *sync.WaitGroup, address, fileName string, stop func()) {
	wg.Add(1)

	go func() {
		err := sendNmea(ctx, wg, address, fileName, stop)
		if err != nil {
			fmt.Printf("ping stopped with error: %s\n", err.Error())
		}
	}()
}

func sendNmea(ctx context.Context, wg *sync.WaitGroup, address, fileName string, stop func()) error {
	defer wg.Done()

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		stop()
		return errors.Wrapf(err, "failed to read file %s", fileName)
	}

	conn, err := net.Dial("udp", address)
	if err != nil {
		stop()
		return errors.Wrap(err, "failed to dial")
	}
	defer conn.Close()

	cnt := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Println("sender is going down")
			stop()
			return nil
		default:
		}

		fmt.Printf("\n\n =========\nSENDING NEXT #%d \n\n", cnt)
		cnt++

		if _, err := conn.Write(data); err != nil {
			fmt.Printf("failed to send data: %s", err.Error())
		}

		time.Sleep(1 * time.Second)
	}
}

//func server(wg *sync.WaitGroup, address, id string) error {
//	defer wg.Done()
//
//	udpAddress, err := net.ResolveUDPAddr("udp", address)
//	if err != nil {
//		return errors.Wrap(err, "server failed to resolve address")
//	}
//
//	conn, err := net.ListenUDP("udp", udpAddress)
//	if err != nil {
//		return errors.Wrap(err, "server failed to listen UDP")
//	}
//	defer func() {
//		_ = conn.Close()
//	}()
//
//	for {
//		buff := make([]byte, 20)
//		rlen, remote, err := conn.ReadFromUDP(buff[:])
//
//		if err != nil {
//			return errors.Wrap(err, "failed to read UDP")
//		}
//
//		message := strings.TrimSpace(string(buff))
//
//		fmt.Printf("server [%s] message from [%s:%d] (%d bytes): %s\n", id, remote.IP, remote.Port, rlen, message)
//	}
//}
