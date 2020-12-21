package exchange

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/adshao/go-binance"

	"github.com/raulk/markdyn/config"
	"github.com/raulk/markdyn/model"
)

// https://github.com/binance/binance-spot-api-docs/blob/master/web-socket-streams.md
type binanceConnector struct {
	sync.Mutex
	config  *config.Binance
	started bool

	stopC  []chan struct{}
	doneC  []chan struct{}
	tradeC chan<- *model.Trade
}

var _ Connector = (*binanceConnector)(nil)

func NewBinanceConnector(cfg *config.Config) Connector {
	return &binanceConnector{config: cfg.Binance}
}

func (c *binanceConnector) ConsumeTrades(into chan<- *model.Trade, symbols ...model.CanonicalSymbol) error {
	c.Lock()
	defer c.Unlock()

	if c.started {
		return fmt.Errorf("already started")
	}

	c.tradeC = into

	symbols = c.config.Mappings.ToExchangeN(symbols...)
	for _, s := range symbols {
		doneC, stopC, err := binance.WsTradeServe(s, c.handleTrade, func(err error) {
			log.Println(err)
		})
		if err != nil {
			return fmt.Errorf("failed to start binance exchange connector")
		}
		c.doneC = append(c.doneC, doneC)
		c.stopC = append(c.stopC, stopC)
	}
	return nil
}

func (c *binanceConnector) handleTrade(evt *binance.WsTradeEvent) {
	symbol := c.config.Mappings.ToCanonical(evt.Symbol)
	price, err := strconv.ParseFloat(evt.Price, 64)
	if err != nil {
		// TODO send elsewhere.
		log.Println(err)
	}
	quantity, err := strconv.ParseFloat(evt.Quantity, 64)
	if err != nil {
		log.Println(err)
	}
	side := model.SideBuyer
	if !evt.IsBuyerMaker {
		side = model.SideSeller
	}
	c.tradeC <- &model.Trade{
		Source:    "binance",
		Timestamp: time.Unix(evt.Time, 0),
		Symbol:    symbol,
		Price:     price,
		Quantity:  quantity,
		Side:      side,
	}
}

func (c *binanceConnector) Close() error {
	for _, c := range c.stopC {
		c <- struct{}{}
	}
	for _, c := range c.doneC {
		<-c
	}
	return nil
}
