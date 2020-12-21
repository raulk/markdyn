package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/raulk/markdyn/config"
	"github.com/raulk/markdyn/exchange"
	"github.com/raulk/markdyn/model"
	"github.com/raulk/markdyn/sink"
)

// Connectors is a registry of all the exchange connectors markdyn knows about.
var Connectors = map[string]func(cfg *config.Config) exchange.Connector{
	"binance":  exchange.NewBinanceConnector,
	"coinbase": exchange.NewCoinbaseConnector,
}

// Sinks is a registry of all sinks markdyn knows about.
var Sinks = map[string]func(cfg *config.Config) sink.Sink{
	"stdout": sink.NewStdout,
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("missing configuration file; usage: markdyn <config_file>")
	}

	cfg := parseConfig(os.Args[1])

	if len(cfg.Sources.Enabled) == 0 {
		log.Fatalf("no exchanges enabled")
	}

	if len(cfg.Sources.Symbols) == 0 {
		log.Fatalf("no symbols enabled")
	}

	exchanges := cfg.Sources.Enabled

	// tradesCh is where connectors will send trades.
	// if there is only one sink, it will consume directly from this channel.
	// if there are more than one, we'll create per-sink channels and a
	// goroutine will multicast each trade to all channels.
	tradesCh := make(chan *model.Trade, 1024)

	// start the sinks.
	sinks := startSinks(cfg, tradesCh)

	// start the exchange connectors.
	connectors := startConnectors(exchanges, cfg, tradesCh)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(
		signalCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)

	<-signalCh

	for _, c := range connectors {
		_ = c.Close()
	}
	close(tradesCh)
	for _, s := range sinks {
		_ = s.Close()
	}

	log.Println("done")

}

func startSinks(cfg *config.Config, tradesCh chan *model.Trade) []sink.Sink {
	sinksChs := []chan *model.Trade{tradesCh}
	if len(cfg.Sinks) > 1 {
		sinksChs = sinksChs[:0]
		for range cfg.Sinks {
			sinksChs = append(sinksChs, make(chan *model.Trade, 1024))
		}
		// multicast goroutine.
		go func() {
			for t := range tradesCh {
				for _, ch := range sinksChs {
					ch <- t
				}
			}
		}()
	}

	var sinks []sink.Sink
	for i, s := range cfg.Sinks {
		ctor, ok := Sinks[s.Kind]
		if !ok {
			log.Fatalf("unrecognized sink: %s", s.Kind)
		}
		snk := ctor(cfg) // construct the sink.
		err := snk.WriteTrades(sinksChs[i])
		if err != nil {
			log.Fatalf("failed to start sink: %s", err)
		}
		sinks = append(sinks, snk)
	}
	return sinks
}

func startConnectors(exchanges []string, cfg *config.Config, ch chan *model.Trade) []exchange.Connector {
	var connectors []exchange.Connector
	for _, id := range exchanges {
		ctor, ok := Connectors[id]
		if !ok {
			log.Fatalf("unrecognized exchange: %s", id)
		}
		e := ctor(cfg) // construct the exchange connector.
		err := e.ConsumeTrades(ch, cfg.Sources.Symbols...)
		if err != nil {
			log.Fatalf("failed to configure exchange: %s", err)
		}
		connectors = append(connectors, e)
	}
	return connectors
}

func parseConfig(path string) *config.Config {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open config file: %s", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read config file: %s", err)
	}

	var cfg config.Config
	if err = json.Unmarshal(bytes, &cfg); err != nil {
		log.Fatalf("failed to parse config file: %s", err)
	}

	return &cfg
}
