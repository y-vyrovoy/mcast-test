package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"

	"sync"

	"github.com/pkg/errors"
)

func main() {

	fmt.Println("msg:", string(generateNmeaMessage(343)))

	//addresses()

	wg := &sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())

	//addr := "224.0.0.1:8888"
	addr := "239.0.112.1:6501"

	//netInterface := "en0"
	//mcast.StartServer(ctx, wg, netInterface, addr, "0", cancel)
	//mcast.StartServer(ctx, wg, netInterface, addr, "1", cancel)

	//mcast.StartPing(ctx, wg, addr, "3", cancel)

	//tctip.StartServer(ctx, wg, ":8080", "0", cancel)
	//tctip.StartServer(ctx, wg, ":8080", "1", cancel)

	//tctip.StartPing(ctx, wg, ":8080", "0", cancel)

	startSendNmea(ctx, wg, addr, "test3.nmea", cancel)

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

	fmt.Println()
	fmt.Println("--------")
	fmt.Println(data)
	fmt.Println(string(data))
	fmt.Println("--------")
	fmt.Println()

	conn, err := net.Dial("udp", address)
	if err != nil {
		stop()
		return errors.Wrap(err, "failed to dial")
	}
	defer conn.Close()

	cnt := 340
	for {
		select {
		case <-ctx.Done():
			fmt.Println("sender is going down")
			stop()
			return nil
		default:
		}

		fmt.Printf("\n\n =========\nSENDING NEXT #%d \n\n", cnt)

		if _, err := conn.Write(generateNmeaMessage(cnt)); err != nil {
			fmt.Printf("failed to send data: %s", err.Error())
		}

		cnt++

		time.Sleep(10 * time.Second)
	}
}

func generateNmeaMessage(i int) []byte {
	body := fmt.Sprintf("HEHDT,%d.0,T", i)

	var sum rune
	for _, b := range body {
		sum = sum ^ b
	}
	checksum := strings.ToUpper(fmt.Sprintf("%2.x", sum))

	return []byte(fmt.Sprintf("$%s*%s\r\n", body, checksum))
}