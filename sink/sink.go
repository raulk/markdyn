package sink

import (
	"io"

	"github.com/raulk/markdyn/model"
)

// Sink is a destination for data collected by markdyn.
type Sink interface {
	io.Closer

	// WriteTrades starts a loop that consumes from ch and writes the incoming
	// trades into a destination (e.g. stdout, file, database).
	//
	// It runs in the background in	a dedicated goroutine.
	WriteTrades(ch <-chan *model.Trade) error
}
