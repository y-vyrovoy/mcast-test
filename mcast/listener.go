package mcast

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

func StartServer(ctx context.Context, wg *sync.WaitGroup, ifName, address, id string, stop func()) {
	wg.Add(1)
	go func() {
		defer stop()

		err := server(ctx, wg, ifName, address, id)
		if err != nil {
			fmt.Printf("server stopped with error: %s\n", err.Error())
		}
	}()
}

func server(ctx context.Context, wg *sync.WaitGroup, ifName, address, id string) error {
	defer wg.Done()

	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return errors.Wrapf(err, "failed to resolve address %s", address)
	}

	netInterface, err := net.InterfaceByName(ifName)
	if err != nil {
		return errors.Wrapf(err, "failed to find %s interface", ifName)
	}

	conn, err := net.ListenMulticastUDP("udp", netInterface, addr)
	if err != nil {
		return errors.Wrap(err, "failed to listen to multicast group")
	}

	buff := make([]byte, 200)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("server is going down")
			return nil
		default:
		}

		rlen, srcAddr, err := conn.ReadFromUDP(buff)
		if err != nil {
			return errors.Wrap(err, "failed to read udp message")
		}

		message := strings.TrimSpace(string(buff))

		fmt.Printf("server [%s] message from [%s:%d] (%d bytes):\n"+
			"%s\n---------\n\n",
			id, srcAddr.IP, srcAddr.Port, rlen, message)
	}
}
