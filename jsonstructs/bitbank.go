package jsonstructs

import "encoding/json"

// BitbankTicker is ticker message from Bitbank.
type BitbankTicker struct {
	Sell      string `json:"sell"`
	Buy       string `json:"buy"`
	High      string `json:"high"`
	Low       string `json:"low"`
	Last      string `json:"last"`
	Volume    string `json:"vol"`
	Timestamp int64  `json:"timestamp"`
}

// BitbankTransactions is transactions message from Bitbank.
type BitbankTransactions struct {
	TransactionID int    `json:"transaction_id"`
	Side          string `json:"side"`
	Amount        string `json:"amount"`
	Price         string `json:"price"`
	ExecutedAt    int64  `json:"excuted_at"`
}

// BitbankDepthDiff is depth diff message from Bitbank.
type BitbankDepthDiff struct {
	Asks      [][]string `json:"a"`
	Bids      [][]string `json:"b"`
	Timestamp int64      `json:"t"`
}

// BitbankDepthWhole is depth whole message from Bitbank.
type BitbankDepthWhole struct {
	Asks      [][]string `json:"asks"`
	Bids      [][]string `json:"bids"`
	Timestamp int64      `json:"timestamp"`
}

// BitbankRoot is root object of Bitbank message.
type BitbankRoot struct {
	RoomName string          `json:"room_name"`
	Message  json.RawMessage `json:"message"`
}

// BitbankWrapper is wrapper for message (possibly socket.io spec).
type BitbankWrapper [2]json.RawMessage

// BitbankSubscribe is subscribe message.
type BitbankSubscribe [2]string

// Initialize initializes subscribe message.
func (w *BitbankSubscribe) Initialize() {
	w[0] = "join-room"
}
