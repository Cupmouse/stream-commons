package jsondef

// BitfinexBook is auto-generated
type BitfinexBook struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Count  int64   `json:"count"`
	Side   string  `json:"side"`
	Size   float64 `json:"size"`
}

// TypeDefBitfinexBook is auto-generated
var TypeDefBitfinexBook = []byte("{\"symbol\": \"symbol\", \"price\": \"price\", \"count\": \"int\", \"side\": \"side\", \"size\": \"size\"}")

// BitfinexTrades is auto-generated
type BitfinexTrades struct {
	Symbol    string  `json:"symbol"`
	OrderID   int64   `json:"orderId"`
	Price     float64 `json:"price"`
	Timestamp string  `json:"timestamp"`
	Side      string  `json:"side"`
	Size      float64 `json:"size"`
}

// TypeDefBitfinexTrades is auto-generated
var TypeDefBitfinexTrades = []byte("{\"symbol\": \"symbol\", \"orderId\": \"int\", \"price\": \"price\", \"timestamp\": \"timestamp\", \"side\": \"side\", \"size\": \"size\"}")
