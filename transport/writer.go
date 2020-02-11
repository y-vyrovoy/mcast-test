package transport

import (
	"fmt"
	"time"

	"main/datasource"
)

type (
	Writer interface {
		Run() error
		Write(data []byte) error
	}

	Sender struct {
		writer Writer
		reader *datasource.MessageReader
	}
)

func NewSender(writer Writer, reader *datasource.MessageReader) *Sender {
	return &Sender{
		writer: writer,
		reader:reader,
	}
}

func (s *Sender) Run() {

	cnt := 0

	for {

		fmt.Println("\n|--->>>>>")

		data, delay, ok := s.reader.ReadNext()

		if !ok {
			fmt.Println("reader ReadNext() returned false")
			return
		}

		fmt.Printf("-- SEND #%d \n", cnt)

		err := s.writer.Write([]byte(data))

		if err != nil {
			fmt.Println("-- !!! ERR: " + err.Error())
			return
		}
		cnt++

		fmt.Println("<<<<---|\n")

		fmt.Printf("Sleep for %v\n\n", delay)
		time.Sleep(delay)
	}
}
