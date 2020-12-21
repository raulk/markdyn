package model

import "encoding/json"

// SymbolMapping maps exchange-specific symbols to canonical mappings, and vice
// versa. The forward mapping is exchange => canonical.
type SymbolMapping struct {
	forward map[string]string
	inverse map[string]string
}

var _ json.Unmarshaler = (*SymbolMapping)(nil)
var _ json.Marshaler = (*SymbolMapping)(nil)

func (sm *SymbolMapping) UnmarshalJSON(bytes []byte) error {
	var m = map[string]string{}
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		return err
	}
	sm.forward = m
	sm.inverse = make(map[string]string, len(m))
	for k, v := range m {
		sm.inverse[v] = k
	}
	return nil
}

func (sm *SymbolMapping) MarshalJSON() ([]byte, error) {
	return json.Marshal(sm.forward)
}

func (sm *SymbolMapping) ToCanonical(symbol ExchangeSymbol) CanonicalSymbol {
	return sm.translate(symbol, sm.forward)
}

func (sm *SymbolMapping) ToExchange(symbol CanonicalSymbol) ExchangeSymbol {
	return sm.translate(symbol, sm.inverse)
}

func (sm *SymbolMapping) ToCanonicalN(symbols ...ExchangeSymbol) []CanonicalSymbol {
	return sm.translateN(symbols, sm.forward)
}

func (sm *SymbolMapping) ToExchangeN(symbols ...CanonicalSymbol) []ExchangeSymbol {
	return sm.translateN(symbols, sm.inverse)
}

func (sm *SymbolMapping) translateN(from []string, table map[string]string) []string {
	ret := make([]string, 0, len(from))
	for _, s := range from {
		ret = append(ret, sm.translate(s, table))
	}
	return ret
}

func (sm *SymbolMapping) translate(from string, table map[string]string) string {
	if other, ok := table[from]; ok {
		return other
	} else {
		return from
	}
}
