package tcp

import (
	"net"

	"github.com/pkg/errors"
)

type (
	connectionWriter struct {
		conn net.Conn
	}
)

func NewConnectionWriter(conn net.Conn) *connectionWriter{
	return &connectionWriter{conn:conn}
}

func (w *connectionWriter) Run() error {
	return nil
}


func (w *connectionWriter) Write(data []byte) error {
	_, err := w.conn.Write(data)

	if err != nil {
		return errors.Wrap(err, "failed to write data")
	}

	return nil
}
