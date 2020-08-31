package jsondef

// BitfinexBook is auto-generated
type BitfinexBook struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Count  int64   `json:"count"`
	Size   float64 `json:"size"`
}

// TypeDefBitfinexBook is auto-generated
var TypeDefBitfinexBook = []byte("{\"symbol\": \"symbol\", \"price\": \"price\", \"count\": \"int\", \"size\": \"size\"}")

// BitfinexTrades is auto-generated
type BitfinexTrades struct {
	Symbol    string  `json:"symbol"`
	OrderID   int64   `json:"orderId"`
	Price     float64 `json:"price"`
	Timestamp string  `json:"timestamp"`
	Size      float64 `json:"size"`
}

// TypeDefBitfinexTrades is auto-generated
var TypeDefBitfinexTrades = []byte("{\"symbol\": \"symbol\", \"orderId\": \"int\", \"price\": \"price\", \"timestamp\": \"timestamp\", \"size\": \"size\"}")
