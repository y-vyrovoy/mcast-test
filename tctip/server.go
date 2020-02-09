package tctip

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/pkg/errors"
	"main/input"
)

func StartServer(ctx context.Context, wg *sync.WaitGroup, address string, inReader *input.MessageReader, stop func()) {
	wg.Add(1)
	go func() {
		defer stop()

		err := server(ctx, wg, address, inReader)
		if err != nil {
			fmt.Printf("server stopped with error: %s\n", err.Error())
		}
	}()
}

func server(ctx context.Context, wg *sync.WaitGroup, address string, inReader *input.MessageReader) error {
	defer wg.Done()

	ipAddress, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return errors.Wrapf(err, "failed to resolve address %s", address)
	}

	l, err := net.ListenTCP("tcp4", ipAddress)
	if err != nil {
		return errors.Wrapf(err, "failed to listen address %s", address)
	}
	defer l.Close()


	for {
		select {
		case <-ctx.Done():
			fmt.Println("server is going down")
			return nil
		default:
		}

		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept connection: %s",  err.Error())
		}

		go msgHandler(conn, *inReader)
	}
}

var cnt = 0

func msgHandler(conn net.Conn, reader input.MessageReader) {

	defer conn.Close()

	for {

		data, delay, ok := reader.ReadNext()
		if !ok {
			return
		}

		fmt.Printf("------ %d SENDING ------- \n%v--------------\n", cnt, data)
		_, err := conn.Write([]byte(data))

		if err != nil {
			fmt.Println("----- ERR: " + err.Error())
			return
		}
		cnt++

		fmt.Printf("sleep for %v\n", delay)
		time.Sleep(delay)
	}
}