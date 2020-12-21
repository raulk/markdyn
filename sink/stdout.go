package sink

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/raulk/markdyn/config"
	"github.com/raulk/markdyn/model"
)

type stdoutSink struct {
	closeCh chan struct{}
	doneCh  chan struct{}
}

var _ Sink = (*stdoutSink)(nil)

func NewStdout(_ *config.Config) Sink {
	return &stdoutSink{
		closeCh: make(chan struct{}),
		doneCh:  make(chan struct{}),
	}
}

func (s *stdoutSink) WriteTrades(ch <-chan *model.Trade) error {
	go func() {
		defer close(s.doneCh)
		defer os.Stdout.Sync()

		for {
			select {
			case t := <-ch:
				bytes, err := json.Marshal(t)
				if err != nil {
					log.Fatalf("failed to marshal json trade: %s", err)
				}
				fmt.Println(string(bytes))

			case <-s.closeCh:
				return
			}
		}
	}()

	return nil
}

func (s *stdoutSink) Close() error {
	close(s.closeCh)
	<-s.doneCh
	return nil
}
