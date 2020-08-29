package jsondef

// LiquidExecutionsCash is auto-generated
type LiquidExecutionsCash struct {
	Pair      string  `json:"pair"`
	CreatedAt int64   `json:"createdAt"`
	ID        int64   `json:"id"`
	Price     float64 `json:"price"`
	Size      float64 `json:"size"`
}

// TypeDefLiquidExecutionsCash is auto-generated
var TypeDefLiquidExecutionsCash = []byte("{\"pair\": \"pair\", \"createdAt\": \"int\", \"id\": \"int\", \"price\": \"price\", \"size\": \"size\"}")

// LiquidPriceLaddersCash is auto-generated
type LiquidPriceLaddersCash struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}

// TypeDefLiquidPriceLaddersCash is auto-generated
var TypeDefLiquidPriceLaddersCash = []byte("{\"price\": \"price\", \"size\": \"size\"}")
