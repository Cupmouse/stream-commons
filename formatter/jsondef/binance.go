package jsondef

// BinanceDepth is auto-generated
type BinanceDepth struct {
	EventTime string  `json:"eventTime"`
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Size      float64 `json:"size"`
}

// TypeDefBinanceDepth is auto-generated
var TypeDefBinanceDepth = []byte("{\"eventTime\": \"timestamp\", \"symbol\": \"symbol\", \"price\": \"price\", \"size\": \"size\"}")

// BinanceRestDepth is auto-generated
type BinanceRestDepth struct {
	Symbol  string  `json:"symbol"`
	OrderID int64   `json:"orderId"`
	Price   float64 `json:"price"`
	Size    float64 `json:"size"`
}

// TypeDefBinanceRestDepth is auto-generated
var TypeDefBinanceRestDepth = []byte("{\"symbol\": \"symbol\", \"orderId\": \"int\", \"price\": \"price\", \"size\": \"size\"}")

// BinanceTrade is auto-generated
type BinanceTrade struct {
	Symbol         string  `json:"symbol"`
	Price          float64 `json:"price"`
	Timestamp      string  `json:"timestamp"`
	Size           float64 `json:"size"`
	TradeID        int64   `json:"tradeID"`
	BuyerOrderID   int64   `json:"buyerOrderID"`
	SellterOrderID int64   `json:"sellterOrderID"`
	EventTime      string  `json:"eventTime"`
}

// TypeDefBinanceTrade is auto-generated
var TypeDefBinanceTrade = []byte("{\"symbol\": \"symbol\", \"price\": \"price\", \"timestamp\": \"timestamp\", \"size\": \"size\", \"tradeID\": \"int\", \"buyerOrderID\": \"int\", \"sellterOrderID\": \"int\", \"eventTime\": \"timestamp\"}")

// BinanceTicker is auto-generated
type BinanceTicker struct {
	EventTime                   string  `json:"eventTime"`
	Symbol                      string  `json:"symbol"`
	PriceChange                 float64 `json:"priceChange"`
	PriceChanePercent           float64 `json:"priceChanePercent"`
	WeightedAveragePrice        float64 `json:"weightedAveragePrice"`
	FirstTradePrice             float64 `json:"firstTradePrice"`
	LastPrice                   float64 `json:"lastPrice"`
	LastQuantity                float64 `json:"lastQuantity"`
	BestBidPrice                float64 `json:"bestBidPrice"`
	BestBidQuantity             float64 `json:"bestBidQuantity"`
	BestAskPrice                float64 `json:"bestAskPrice"`
	BestAskQuantity             float64 `json:"bestAskQuantity"`
	OpenPrice                   float64 `json:"openPrice"`
	HighPrice                   float64 `json:"highPrice"`
	LowPrice                    float64 `json:"lowPrice"`
	TotalTradedBaseAssetVolume  float64 `json:"totalTradedBaseAssetVolume"`
	TotalTradedQuoteAssetVolume float64 `json:"totalTradedQuoteAssetVolume"`
	StatisticsOpenTime          string  `json:"statisticsOpenTime"`
	StatisticsCloseTime         string  `json:"statisticsCloseTime"`
	FirstTradeID                int64   `json:"firstTradeID"`
	LastTradeID                 int64   `json:"lastTradeID"`
	TotalNumberOfTrades         int64   `json:"totalNumberOfTrades"`
}

// TypeDefBinanceTicker is auto-generated
var TypeDefBinanceTicker = []byte("{\"eventTime\": \"timestamp\", \"symbol\": \"symbol\", \"priceChange\": \"float\", \"priceChanePercent\": \"float\", \"weightedAveragePrice\": \"float\", \"firstTradePrice\": \"float\", \"lastPrice\": \"float\", \"lastQuantity\": \"float\", \"bestBidPrice\": \"float\", \"bestBidQuantity\": \"float\", \"bestAskPrice\": \"float\", \"bestAskQuantity\": \"float\", \"openPrice\": \"float\", \"highPrice\": \"float\", \"lowPrice\": \"float\", \"totalTradedBaseAssetVolume\": \"float\", \"totalTradedQuoteAssetVolume\": \"float\", \"statisticsOpenTime\": \"timestamp\", \"statisticsCloseTime\": \"timestamp\", \"firstTradeID\": \"int\", \"lastTradeID\": \"int\", \"totalNumberOfTrades\": \"int\"}")
