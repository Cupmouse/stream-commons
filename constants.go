package streamcommons

// Name of exchanges
const (
	ExchangeBitflyer = "bitflyer"
	ExchangeBitfinex = "bitfinex"
	ExchangeBitmex   = "bitmex"
	ExchangeBinance  = "binance"
	ExchangeLiquid   = "liquid"
)

// Bitfinex related
const (
	BitfinexChannelPrefixBook   = "book_"
	BitfinexChannelPrefixTrades = "trades_"
)

// Bitflyer related
const (
	BitflyerChannelPrefixLightningExecutions    = "lightning_executions_"
	BitflyerChannelPrefixLightningBoard         = "lightning_board_"
	BitflyerchannelPrefixLightningBoardSnapshot = "lightning_board_snapshot_"
	BitflyerchannelPrefixLightningTicker        = "lightning_ticker_"
)

// Bitmex related
const (
	BitmexChannelOrderBookL2 = "orderBookL2"
	BitmexChannelTrade       = "trade"
	BitmexChannelInstrument  = "instrument"
	BitmexChannelLiquidation = "liquidation"
	BitmexChannelSettlement  = "settlement"
	BitmexChannelInsurance   = "insurance"
	BitmexChannelFunding     = "funding"
)

// Binance related
const (
	// BinanceStreamRESTDepth is the prefix for binance indicating a REST depth channel.
	BinanceStreamRESTDepth   = "rest_depth"
	BinanceStreamDepth       = "depth@100ms"
	BinanceStreamTrade       = "trade"
	BinanceStreamTicker      = "ticker"
	BinancePricePrecision    = 8
	BinanceQuantityPrecision = 8
)

// Liquid related
const (
	LiquidChannelConnectionEstablished = "connection_established"
	LiquidChannelPrefixLaddersCash     = "price_ladders_cash_"
	LiquidChannelPrefixExecutionsCash  = "executions_cash_"
)

// Common format
const (
	CommonFormatSell    = "Sell"
	CommonFormatBuy     = "Buy"
	CommonFormatUnknown = "Unknown"
)

// StateChannelSubscribed is the channel name for subscribed channel message in status line
const StateChannelSubscribed = "!subscribed"

// ChannelUnknown is the placeholder for message whose channel could not be specified
const ChannelUnknown = "!unknown"
