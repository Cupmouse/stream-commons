package jsondef

// LiquidExecutionsCash is auto-generated
type LiquidExecutionsCash struct {
	Symbol    string  `json:"symbol"`
	CreatedAt string  `json:"createdAt"`
	ID        int64   `json:"id"`
	Price     float64 `json:"price"`
	Side      string  `json:"side"`
	Size      float64 `json:"size"`
}

// TypeDefLiquidExecutionsCash is auto-generated
var TypeDefLiquidExecutionsCash = []byte("{\"symbol\": \"symbol\", \"createdAt\": \"timestamp\", \"id\": \"int\", \"price\": \"price\", \"side\": \"side\", \"size\": \"size\"}")

// LiquidPriceLaddersCash is auto-generated
type LiquidPriceLaddersCash struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Side   string  `json:"side"`
	Size   float64 `json:"size"`
}

// TypeDefLiquidPriceLaddersCash is auto-generated
var TypeDefLiquidPriceLaddersCash = []byte("{\"symbol\": \"symbol\", \"price\": \"price\", \"side\": \"side\", \"size\": \"size\"}")
