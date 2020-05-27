package jsonstructs

// BitfinexBook is the root structure of orderbook data from bitfinex public orderbook channel
type BitfinexBook = [2]interface{}

// BitfinexBookOrder is single order in the orderbook from bitfinex public channel
// price, count, amount
type BitfinexBookOrder = []interface{}

// BitfinexTrades is the root array of trades data from bitfinex
type BitfinexTrades = []interface{}

// BitfinexTradesElement is the element of bitfinex trades
type BitfinexTradesElement = []interface{}

// BitfinexStatusSubscribed provides structure for subscribed line in status line.
// map[channel]chanID
type BitfinexStatusSubscribed = map[string]int

// BitfinexSubscribe is the root struct of subscribe message sent to server
type BitfinexSubscribe struct {
	Event   string  `json:"event"`
	Channel string  `json:"channel"`
	Symbol  string  `json:"symbol"`
	Prec    *string `json:"prec"`
	Frec    *string `json:"frec"`
	Len     *string `json:"len"`
}

// Initialize will initialize event field to "subscribe" if this struct was created in go
func (s *BitfinexSubscribe) Initialize() {
	s.Event = "subscribe"
}

// BitfinexSubscribed is the special message from bitfinex either it is the response to subscribe.
type BitfinexSubscribed struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	ChanID  int    `json:"chanId"`
	Symbol  string `json:"symbol"`
	Pair    string `json:"pair"`
}

// Initialize will initialize event field to "subscribed" if this struct was created in go
func (s *BitfinexSubscribed) Initialize() {
	s.Event = "subscribed"
}
