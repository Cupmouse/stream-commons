package jsondef

// BitfinexBook is auto-generated
type BitfinexBook struct {
	Pair  string  `json:"pair"`
	Price float64 `json:"price"`
	Count int64   `json:"count"`
	Size  float64 `json:"size"`
}

// TypeDefBitfinexBook is auto-generated
var TypeDefBitfinexBook = []byte("{\"pair\": \"pair\", \"price\": \"price\", \"count\": \"int\", \"size\": \"size\"}")

// BitfinexTrades is auto-generated
type BitfinexTrades struct {
	Pair      string  `json:"pair"`
	OrderID   int64   `json:"orderId"`
	Price     float64 `json:"price"`
	Timestamp string  `json:"timestamp"`
	Size      float64 `json:"size"`
}

// TypeDefBitfinexTrades is auto-generated
var TypeDefBitfinexTrades = []byte("{\"pair\": \"pair\", \"orderId\": \"int\", \"price\": \"price\", \"timestamp\": \"timestamp\", \"size\": \"size\"}")
