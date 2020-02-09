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

func StartClient(ctx context.Context, wg *sync.WaitGroup, address, id string, inReader *input.MessageReader, stop func()) {
	wg.Add(1)

	go func() {
		defer stop()

		err := ping(ctx, wg, address, id, inReader)
		if err != nil {
			fmt.Printf("ping stopped with error: %s\n", err.Error())
		}
	}()
}

func ping(ctx context.Context, wg *sync.WaitGroup, address, id string, reader *input.MessageReader) error {
	defer wg.Done()

	fmt.Println("let's send")

	cnt := 0

	//rAddr, err := net.ResolveTCPAddr("tcp", address)
	//if err != nil {
	//	panic("can't resolve remote address: " + err.Error())
	//}
	//
	//lAddr, err := net.ResolveTCPAddr("tcp", "192.168.15.137:8081")
	//if err != nil {
	//	panic("can't resolve remote address: " + err.Error())
	//}
	//
	//
	//fmt.Println("addresses resolved")
	//
	//conn, err := net.DialTCP("tcp", lAddr, rAddr)

	conn, err := net.Dial("tcp", address)

	fmt.Println("dialed")

	if err != nil {
		return errors.Wrap(err, "failed to dial")
	}
	defer conn.Close()

	fmt.Println("start sending")

	for {

		data, delay, ok := reader.ReadNext()
		if !ok {
			return nil
		}

		fmt.Printf("------ %d SENDING ------- \n%v--------------\n", cnt, data)
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

