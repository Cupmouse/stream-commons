package jsonstructs

import "encoding/json"

// BitflyerTickerParamsMessage is the ticker
type BitflyerTickerParamsMessage struct {
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

// BitflyerExecutionsParamMessageElement is the executed order element in executions
type BitflyerExecutionsParamMessageElement struct {
	ID                         uint64  `json:"id"`
	Side                       string  `json:"side"`
	Price                      float64 `json:"price"`
	Size                       float64 `json:"size"`
	ExecDate                   string  `json:"exec_date"`
	BuyChildOrderAcceptanceID  string  `json:"buy_child_order_acceptance_id"`
	SellChildOrderAcceptanceID string  `json:"sell_child_order_acceptance_id"`
}

// BitflyerBoardParamsMessageOrder is the order of orderbook in bitflyer
type BitflyerBoardParamsMessageOrder struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}

// BitflyerBoardParamsMessage is actual payload for orderbook from bitflyer
type BitflyerBoardParamsMessage struct {
	Asks []BitflyerBoardParamsMessageOrder `json:"asks"`
	Bids []BitflyerBoardParamsMessageOrder `json:"bids"`
}

// BitflyerRoot is the root of message from bitflyer
type BitflyerRoot struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Channel string          `json:"channel"`
		Message json.RawMessage `json:"message"`
	} `json:"params"`
}

// Initialize initialize constants on this struct
func (b *BitflyerRoot) Initialize() {
	b.JSONRPC = "2.0"
	b.Method = "channelMessage"
}

// BitflyerSubscribe is the root of subscribe message client sends to server
type BitflyerSubscribe struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  struct {
		Channel string `json:"channel"`
	} `json:"params"`
}

// Initialize initialize constants on this struct
func (b *BitflyerSubscribe) Initialize() {
	b.JSONRPC = "2.0"
	b.Method = "subscribe"
}

// BitflyerSubscribed is the root of orderbook message from bitflyer
type BitflyerSubscribed struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  bool   `json:"result"`
}

// Initialize initialize constants on this struct
func (b *BitflyerSubscribed) Initialize() {
	b.JSONRPC = "2.0"
}

// BitflyerStateSubscribed is a list of subscribed channels listed in state line in dataset
type BitflyerStateSubscribed []string
