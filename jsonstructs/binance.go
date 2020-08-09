package jsonstructs

import "encoding/json"

// BinanceRateLimiter is the rate limit information used in REST API
type BinanceRateLimiter struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
}

// BinanceFilter is the filter used in REST API
// type BinanceFilter struct {
// 	FilterType string `json:"filterType"`
// 	MinPrice   string `json:"minPrice"`
// 	MaxPrice   string `json:"maxPrice"`
// 	TickSize   string `json:"tickSize"`
// }

// BinanceExchangeInfoSymbol is the individual information for a symbol in BinanceExchangeInfo
type BinanceExchangeInfoSymbol struct {
	Symbol                 string        `json:"symbol"`
	Status                 string        `json:"status"`
	BaseAsset              string        `json:"baseAsset"`
	BaseAssetPrecision     int           `json:"baseAssetPrecision"`
	QuoteAsset             string        `json:"quoteAsset"`
	QuotePrecision         int           `json:"quotePrecision"`
	QuoteAssetPrecision    int           `json:"quoteAssetPrecision"`
	OrderTypes             []string      `json:"orderTypes"`
	IcebergAllowed         bool          `json:"icebergAllowed"`
	OCOAllowed             bool          `json:"ocoAllowed"`
	IsSpoTradingAllowed    bool          `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed bool          `json:"isMarginTradingAlloed"`
	Filter                 []interface{} `json:"filters"`
	Permissons             []string      `json:"permissions"`
}

// BinanceExchangeInfo is the response from exchangeInfo REST API Endpoint
type BinanceExchangeInfo struct {
	Timezone        string                      `json:"timezone"`
	ServerTime      int64                       `json:"serverTime"`
	RateLimits      []BinanceRateLimiter        `json:"rateLimits"`
	ExchangeFilters []interface{}               `json:"exchangeFilters"`
	Symbols         []BinanceExchangeInfoSymbol `json:"symbols"`
}

// BinanceTicker24HrElement is the individual ticker element from REST API
type BinanceTicker24HrElement struct {
	Symbol             string `json:"symbol"`
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	WeightedAvgPrice   string `json:"weightedAvgPrice"`
	PrevClosePrice     string `json:"prevClosePrice"`
	LastPrice          string `json:"lastPrice"`
	LastQty            string `json:"lastQty"`
	BidPrice           string `json:"bidPrice"`
	AskPrice           string `json:"askPrice"`
	OpenPrice          string `json:"openPrice"`
	HighPrice          string `json:"highPrice"`
	LowPrice           string `json:"lowPrice"`
	Volume             string `json:"volume"`
	QuoteVolume        string `json:"quoteVolume"`
	OpenTime           int64  `json:"openTime"`
	CloseTime          int64  `json:"closeTime"`
	FirstID            int64  `json:"firstId"`
	LastID             int64  `json:"lastId"`
	Count              int64  `json:"count"`
}

// BinanceTrade is the individual executed order
type BinanceTrade struct {
	EventType          string `json:"e"`
	EventTime          int64  `json:"E"`
	Symbol             string `json:"s"`
	TradeID            int64  `json:"t"`
	Price              string `json:"p"`
	Quantity           string `json:"q"`
	BuyerOrderID       int64  `json:"b"`
	SellerOrderID      int64  `json:"a"`
	TradeTime          int64  `json:"T"`
	IsBuyerMarketMaker bool   `json:"m"`
	M                  bool   `json:"M"`
}

// BinanceDepthStream is the message of depth stream
type BinanceDepthStream struct {
	// Constant depthUpdate
	DepthUpdate   string `json:"e"`
	EventTime     int64  `json:"E"`
	Symbol        string `json:"s"`
	FirstUpdateID int64  `json:"U"`
	FinalUpdateID int64  `json:"u"`
	// BinanceDepthOrder is the individual order in depth response
	// [0] = price, [1] = quantity
	Bids [][]string `json:"b"`
	Asks [][]string `json:"a"`
}

// BinanceDepthREST is the response of depth REST
type BinanceDepthREST struct {
	LastUpdateID int64 `json:"lastUpdateId"`
	// [0] = price, [1] = quantity
	Bids [][]string `json:"bids"`
	Asks [][]string `json:"asks"`
}

// BinanceSubscribe is for subscribe requests
type BinanceSubscribe struct {
	// Constant SUBSCRIBE
	Method string `json:"method"`
	// Channels to subscribe
	Params []string `json:"params"`
	ID     int      `json:"id"`
}

// Initialize initialize this BinanceSubscribe
// It sets `Method` field to "SUBSCRIBE"
func (s *BinanceSubscribe) Initialize() {
	s.Method = "SUBSCRIBE"
}

// BinanceSubscribeResponse is the response for subscribe request
type BinanceSubscribeResponse struct {
	// Nil if successful
	Result interface{} `json:"result"`
	ID     int         `json:"id"`
}

// BinanceReponseRoot is the root structure for stream messages
type BinanceReponseRoot struct {
	// Name of channel
	Stream string `json:"stream"`
	// Payload
	Data json.RawMessage `json:"data"`
}
