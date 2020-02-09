package mcast

import (
	"context"
	"fmt"
	"main/input"
	"net"
	"sync"
	"time"

	"github.com/pkg/errors"
)

func StartPing(ctx context.Context, wg *sync.WaitGroup, address string, reader *input.MessageReader, stop context.CancelFunc) {
	wg.Add(1)

	go func() {
		defer stop()

		err := ping(ctx, wg, address, reader)
		if err != nil {
			fmt.Printf("ping stopped with error: %s\n", err.Error())
		}
	}()
}

func ping(ctx context.Context, wg *sync.WaitGroup, address string, reader *input.MessageReader) error {
	defer wg.Done()

	cnt := 0

	conn, err := net.Dial("udp", address)
	if err != nil {
		return errors.Wrap(err, "failed to dial")
	}
	defer conn.Close()

	fmt.Println("start sending")

	for {

		data, delay, ok := reader.ReadNext()
		if !ok {
			fmt.Println("reader ReadNext() returned false")
			return nil
		}

		fmt.Printf("------ %d SENDING ------- \n %v--------------\n", cnt, data)
		_, err = conn.Write([]byte(data))

		if err != nil {
			fmt.Println("----- ERR: " + err.Error())
			return err
		}
		//_, _ = fmt.Fprintf(conn, "ping [id:%s] TCP message %d", id, cnt)
		cnt++

		time.Sleep(delay)
	}
}
