package streamcommons

// Binance related
const (
	// BinanceStreamRESTDepth is the prefix for binance indicating a REST depth channel.
	BinanceStreamRESTDepth   = "rest_depth"
	BinancePricePrecision    = 8
	BinanceQuantityPrecision = 8
)

// StateChannelSubscribed is the channel name for subscribed channel message in status line
const StateChannelSubscribed = "!subscribed"

// ChannelUnknown is the placeholder for message whose channel could not be specified
const ChannelUnknown = "!unknown"
