package tctip

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

func StartServer(ctx context.Context, wg *sync.WaitGroup, address, id string, stop func()) {
	wg.Add(1)
	go func() {
		defer stop()

		err := server(ctx, wg, address, id)
		if err != nil {
			fmt.Printf("server stopped with error: %s\n", err.Error())
		}
	}()
}

func server(ctx context.Context, wg *sync.WaitGroup, address, id string) error {
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

	f, _ := l.File()
	fd := f.Fd()

	err = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR,  1)
	if err != nil {
		return errors.Wrap(err, "failed to set SO_REUSEADDR")
	}

	err = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEPORT,  1)
	if err != nil {
		return errors.Wrap(err, "failed to set SO_REUSEPORT")
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Println("server is going down")
			return nil
		default:
		}

		conn, err := l.Accept()
		if err != nil {
			fmt.Errorf("[%s] failed to accept connection: %s", id, err.Error())
		}

		go msgHandler(conn, id)
	}
}

var cnt = 0

func msgHandler(conn net.Conn, id string) {

	defer conn.Close()

	//buff := make([]byte, 200)

	for {
	//	rlen, err := conn.Read(buff)
	//	if err != nil {
	//		fmt.Errorf("[%s] failed to read data: %s", id, err.Error())
	//		return
	//	}
	//
	//	fmt.Printf("[%s] message:\n"+
	//		"\tlen: %d\n"+
	//		"\tsrc: %s\n" +
	//		"\tmsg: %s\n"+
	//		"-------\n\n",
	//		id, rlen, conn.LocalAddr(), string(buff))

		msg := generateNmeaMessage(cnt)

		_, _ = conn.Write(msg)
		fmt.Printf( "------ %d SENT ------- \n", cnt)
		cnt++

		time.Sleep(1 * time.Second)
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