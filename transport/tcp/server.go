package tcp

import (
	"fmt"
	"net"
	"sync"

	"main/sender"
)

type (
	tcpServerWriter struct {
		address     string
		listener    net.Listener
		onConn      onConnectionCB
	}
)

type onConnectionCB func(w sender.Writer)

func NewServerWriter(address string, onConn onConnectionCB) *tcpServerWriter {
	return &tcpServerWriter{
		address:     address,
		onConn:      onConn,
	}
}

func (w *tcpServerWriter) Run(wg *sync.WaitGroup) error {
	conn, err := net.Listen("tcp", w.address)
	if err != nil {
		return err
	}

	w.listener = conn

	wg.Add(1)
	go func() {
		defer wg.Done()
		w.listenWorker()
	}()

	return nil
}

func (w *tcpServerWriter) listenWorker() {

	for {
		conn, err := w.listener.Accept()

		if err != nil {
			fmt.Printf("failed to accept connection from %s to %s: %s",
				conn.RemoteAddr().String(),
				conn.LocalAddr().String(),
				err.Error())

			return
		}

		fmt.Printf("\n ---> connection from %s to %s\n\n",
			conn.RemoteAddr().String(),
			conn.LocalAddr().String())

		go func() {
			w.connectionWorker(conn)
		}()
	}
}

func (w *tcpServerWriter) connectionWorker(conn net.Conn) {

	connWriter := NewConnectionWriter(conn)

	w.onConn(connWriter)

	fmt.Printf("\n---> all data is sent to [%s] \n", conn.RemoteAddr())
	_ = conn.Close()
}
