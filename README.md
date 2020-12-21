# markdyn

markdyn is a tool to study crypto exchange market dynamics. It starts life being
a data stream aggregator. It may do more in the future.

Right now, I'm focusing on getting a normalised, consolidated data stream for:

* trades
* L2 order books
* L3 order flow (for exchanges that support it)

## Implementation details

Connection to exchanges is done via adapters. Adapters need to be configured
with their respective auth tokens. Adapters are responsible for transforming
events to the canonical data model.

Asset and pair selection to watch is done via configuration. Since pairs
symbols can be exchange-specific, users can specify mappings to canonical pairs.

Right now, markdyn supports writing the event stream to an ndjson sink. In the
future, other sinks (e.g. databases) will be supported.
