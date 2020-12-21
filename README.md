# markdyn

markdyn is a tool to study crypto exchange market dynamics. It starts life being
a data stream aggregator. It may do more in the future.

Right now, I'm focusing on getting a normalised, consolidated data stream for
the following items, in descending order of priority:

1. trades
2. L2 order books
3. L3 order flow (for exchanges that support it)

## Exchanges supported

- ✅ Coinbase
- 🚧 Binance
- 🚧 Huobi
- 🚧 Gate.io
- (add your own; submit a PR!)

## Sinks supported

- ✅ stdout (ndjson)
- 🚧 rotating files (ndjson)
- 🚧 databases -- TBD.
- (add your own; submit a PR!)

## Try it out

```shell
$ git clone https://github.com/raulk/markdyn.git
$ cd markdyn
$ go build .
$ ./markdyn example-config.json
```

Modify th configuration to track more assets.

## Implementation details

Connection to exchanges is done via connectors. Connectors need to be configured
with their respective auth tokens. Connectors are responsible for transforming
incoming events to the canonical data model.

Asset and pair selection to watch is done via configuration. Since symbols ids
can be exchange-specific, users can specify mappings to canonical symbols.

The consolidated data stream is then written to one or many sinks. The following
sinks are under development:

1. simple ndjson on stdout.
2. ndjson on rotating log files.
3. databases.

Currently, markdyn does not perform relative ordering of events across sources,
but it may do so through a buffering mechanism in the future.
