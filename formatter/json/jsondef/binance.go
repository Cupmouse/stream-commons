package jsondef

// BinanceDepth is auto-generated
type BinanceDepth struct {
	Pair  string  `json:"pair"`
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}

// TypeDefBinanceDepth is auto-generated
var TypeDefBinanceDepth = []byte("{\"pair\": \"pair\", \"price\": \"price\", \"size\": \"size\"}")

// BinanceRestDepth is auto-generated
type BinanceRestDepth struct {
	Pair    string  `json:"pair"`
	OrderID int64   `json:"orderId"`
	Price   float64 `json:"price"`
	Size    float64 `json:"size"`
}

// TypeDefBinanceRestDepth is auto-generated
var TypeDefBinanceRestDepth = []byte("{\"pair\": \"pair\", \"orderId\": \"int\", \"price\": \"price\", \"size\": \"size\"}")

// BinanceTrade is auto-generated
type BinanceTrade struct {
	Pair           string  `json:"pair"`
	Price          float64 `json:"price"`
	Timestamp      string  `json:"timestamp"`
	Size           float64 `json:"size"`
	TradeID        int64   `json:"tradeID"`
	BuyerOrderID   int64   `json:"buyerOrderID"`
	SellterOrderID int64   `json:"sellterOrderID"`
	Eventtime      string  `json:"eventtime"`
	IsBuyerMaker   bool    `json:"isBuyerMaker"`
}

// TypeDefBinanceTrade is auto-generated
var TypeDefBinanceTrade = []byte("{\"pair\": \"pair\", \"price\": \"price\", \"timestamp\": \"timestamp\", \"size\": \"size\", \"tradeID\": \"int\", \"buyerOrderID\": \"int\", \"sellterOrderID\": \"int\", \"eventtime\": \"timestamp\", \"isBuyerMaker\": \"boolean\"}")
