package tctip

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/pkg/errors"
)

func StartPing(ctx context.Context, wg *sync.WaitGroup, address, id string, stop func()) {
	wg.Add(1)

	go func() {
		defer stop()

		err := ping(ctx, wg, address, id)
		if err != nil {
			fmt.Printf("ping stopped with error: %s\n", err.Error())
		}
	}()
}

func ping(ctx context.Context, wg *sync.WaitGroup, address, id string) error {
	defer wg.Done()

	cnt := 0

	for {
		fmt.Println("let's send")

		conn, err := net.Dial("tcp", address)

		if err != nil {
			return errors.Wrap(err, "failed to dial")
		}
		defer conn.Close()


		ticker := time.NewTicker(1 * time.Second)
		select {
		case <-ctx.Done():
			fmt.Println("sender is going down")
			return nil
		case <-ticker.C:
		}

		_, _ = fmt.Fprintf(conn, "ping [id:%s] TCP message %d", id, cnt)
		cnt++
	}

}