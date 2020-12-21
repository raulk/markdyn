package exchange

import (
	"io"

	"github.com/raulk/markdyn/model"
)

// Connector is an exchange connector. Its role is to obtain real-time data
// from an exchange, and transform it to our canonical data model.
//
// Currently it supports trades, but it will be extended to support L2 and
// L3 order book data.
type Connector interface {
	io.Closer

	// ConsumeTrades starts streaming trade data from this exchange for the
	// specified symbols into the supplied channel.
	//
	// It runs in the background in	a dedicated goroutine.
	ConsumeTrades(into chan<- *model.Trade, symbols ...model.CanonicalSymbol) error
}
