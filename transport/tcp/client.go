package tcp

import (
	"context"
	"net"
	"sync"

	"github.com/pkg/errors"
)

type tcpClientWriter struct {
	address string
	conn    net.Conn

	ctx  context.Context
	wg   *sync.WaitGroup
	stop context.CancelFunc
}

func NewClientWriter(address string) *tcpClientWriter {
	return &tcpClientWriter{
		address: address,
	}
}

func (w *tcpClientWriter) Run() error {

	conn, err := net.Dial("tcp", w.address)

	if err != nil {
		return errors.Wrap(err, "failed to run TCP client sender")
	}
	w.conn = conn

	return nil
}

func (w *tcpClientWriter) Write(data []byte) error {

	_, err := w.conn.Write(data)

	if err != nil {
		return errors.Wrap(err, "failed to send TCP packet")
	}

	return nil
}
