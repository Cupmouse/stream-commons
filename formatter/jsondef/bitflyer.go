package jsondef

// BitflyerBoard is auto-generated
type BitflyerBoard struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Side   string  `json:"side"`
	Size   float64 `json:"size"`
}

// TypeDefBitflyerBoard is auto-generated
var TypeDefBitflyerBoard = []byte("{\"symbol\": \"symbol\", \"price\": \"price\", \"side\": \"side\", \"size\": \"size\"}")

// BitflyerExecutions is auto-generated
type BitflyerExecutions struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Side   string  `json:"side"`
	Size   float64 `json:"size"`
}

// TypeDefBitflyerExecutions is auto-generated
var TypeDefBitflyerExecutions = []byte("{\"symbol\": \"symbol\", \"price\": \"price\", \"side\": \"side\", \"size\": \"size\"}")

// BitflyerTicker is auto-generated
type BitflyerTicker struct {
	ProductCode     string  `json:"product_code"`
	Timestamp       string  `json:"timestamp"`
	TickID          int64   `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	Ltp             float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}

// TypeDefBitflyerTicker is auto-generated
var TypeDefBitflyerTicker = []byte("{\"product_code\": \"string\", \"timestamp\": \"string\", \"tick_id\": \"int\", \"best_bid\": \"float\", \"best_ask\": \"float\", \"best_bid_size\": \"float\", \"best_ask_size\": \"float\", \"total_bid_depth\": \"float\", \"total_ask_depth\": \"float\", \"ltp\": \"float\", \"volume\": \"float\", \"volume_by_product\": \"float\"}")
