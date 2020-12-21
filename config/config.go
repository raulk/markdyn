package config

import "github.com/raulk/markdyn/model"

type Config struct {
	Sources  Sources   `json:"sources"`
	Sinks    []Sink    `json:"sinks"`
	Coinbase *Coinbase `json:"coinbase"`
	Binance  *Binance  `json:"binance"`
}

type Sources struct {
	Enabled []string                `json:"enabled"`
	Symbols []model.CanonicalSymbol `json:"symbols"`
}

type Sink struct {
	Kind string `json:"kind"`
}

type Coinbase struct {
	Mappings *model.SymbolMapping `json:"mappings"`
}

type Binance struct {
	AuthKey   string               `json:"auth_key"`
	SecretKey string               `json:"secret_key"`
	Mappings  *model.SymbolMapping `json:"mappings"`
}
