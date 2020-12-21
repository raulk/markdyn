package model

import "time"

// CanonicalSymbol represents a canonical symbol.
type CanonicalSymbol = string

// ExchangeSymbol represents an exchange symbol.
type ExchangeSymbol = string

// Side indicates who provided liquidity for this trade to execute. If it was
// the buyer, this was a downtick. If it was the seller, this was an uptick.
type Side string

const (
	SideBuyer  = Side("b")
	SideSeller = Side("s")
)

// Trade represents a trade reported by an exchange.
type Trade struct {
	Source    string
	Timestamp time.Time
	Symbol    CanonicalSymbol
	Price     float64
	Quantity  float64
	Side      Side
}
