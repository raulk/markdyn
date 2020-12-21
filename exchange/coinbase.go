package exchange

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/raulk/markdyn/config"
	"github.com/raulk/markdyn/model"

	"github.com/preichenberger/go-coinbasepro/v2"
)

const CoinbaseWssEndpoint = "wss://ws-feed.pro.coinbase.com"

type coinbaseConnector struct {
	sync.Mutex
	config  *config.Coinbase
	started bool

	wsConn *websocket.Conn
	tradeC chan<- *model.Trade
	doneCh chan struct{}
}

var _ Connector = (*coinbaseConnector)(nil)

func NewCoinbaseConnector(cfg *config.Config) Connector {
	return &coinbaseConnector{
		config: cfg.Coinbase,
		doneCh: make(chan struct{}),
	}
}

func (c *coinbaseConnector) ConsumeTrades(into chan<- *model.Trade, symbols ...model.CanonicalSymbol) error {
	c.Lock()
	defer c.Unlock()

	if c.started {
		return fmt.Errorf("already started")
	}

	c.tradeC = into

	var err error
	c.wsConn, _, err = websocket.DefaultDialer.Dial(CoinbaseWssEndpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to dial coinbase websocket endpoint: %w", err)
	}

	symbols = c.config.Mappings.ToExchangeN(symbols...)
	sub := coinbasepro.Message{
		Type: "subscribe",
		Channels: []coinbasepro.MessageChannel{{
			Name:       "heartbeat",
			ProductIds: symbols,
		}, {
			Name:       "ticker",
			ProductIds: symbols,
		}},
	}
	if err = c.wsConn.WriteJSON(sub); err != nil {
		return fmt.Errorf("failed to subscribe to coinbase channels: %w", err)
	}

	go c.consume()

	return nil
}

func (c *coinbaseConnector) consume() {
	defer close(c.doneCh)

	var msg coinbasepro.Message
	for {
		if err := c.wsConn.ReadJSON(&msg); err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				return // we have been closed, so silence this error.
			}
			log.Fatal(err)
		}
		if msg.Type != "ticker" {
			continue
		}
		price, err := strconv.ParseFloat(msg.Price, 64)
		if err != nil {
			log.Println(err)
		}
		quantity, err := strconv.ParseFloat(msg.LastSize, 64)
		if err != nil {
			log.Println(err)
		}
		side := model.SideBuyer
		if msg.Side == "sell" {
			side = model.SideSeller
		}
		c.tradeC <- &model.Trade{
			Source:    "coinbase",
			Timestamp: msg.Time.Time(),
			Symbol:    c.config.Mappings.ToCanonical(msg.ProductID),
			Price:     price,
			Quantity:  quantity,
			Side:      side,
		}
	}
}

func (c *coinbaseConnector) Close() error {
	c.Lock()
	c.started = false
	c.Unlock()

	if err := c.wsConn.SetReadDeadline(time.Now().Add(50 * time.Millisecond)); err != nil {
		return fmt.Errorf("failed to set deadline on underlying WS connection: %w", err)
	}

	<-c.doneCh
	return nil
}
