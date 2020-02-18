package mcast

import (
	"context"
	"net"
	"sync"

	"github.com/pkg/errors"
)

type mcastWriter struct {
	address string
	conn    net.Conn

	ctx  context.Context
	wg   *sync.WaitGroup
	stop context.CancelFunc
}

func NewWriter(address string) *mcastWriter {
	return &mcastWriter{
		address: address,
	}
}

func (w *mcastWriter) Run() error {

	conn, err := net.Dial("udp", w.address)

	if err != nil {
		return errors.Wrap(err, "failed to run UDP sender")
	}
	w.conn = conn

	return nil
}

func (w *mcastWriter) Write(data []byte) error {

	_, err := w.conn.Write(data)

	if err != nil {
		return errors.Wrap(err, "failed to send UDP packet")
	}

	return nil
}
