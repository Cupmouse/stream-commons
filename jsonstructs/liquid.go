package jsonstructs

import "encoding/json"

// LiquidProduct is individual product on Liquid
type LiquidProduct struct {
	ID                  string   `json:"id"`
	ProductType         string   `json:"product_type"`
	Code                string   `json:"code"`
	Name                *string  `json:"name"`
	MarketAsk           *float64 `json:"market_ask"`
	MarketBid           *float64 `json:"market_bid"`
	Indicator           *int     `json:"indicator"`
	Currency            string   `json:"currency"`
	CurrencyPairCode    string   `json:"currency_pair_code"`
	Symbol              *string  `json:"symbol"`
	BtcMinimumWithdraw  *string  `json:"btc_minimum_withdraw"`
	FiatMinimumWithdraw *string  `json:"fiat_minimum_withdraw"`
	PusherChannel       string   `json:"pusher_channel"`
	TakerFee            string   `json:"taker_fee"`
	MakerFee            string   `json:"maker_fee"`
	LowMarketBid        string   `json:"low_market_bid"`
	HighMarketAsk       string   `json:"high_market_ask"`
	Volume24h           *string  `json:"volume_24h"`
	LastPrice24h        *string  `json:"last_price_24h"`
	LastTradedPrice     *string  `json:"last_traded_price"`
	LastTradedQuantity  *string  `json:"last_traded_quantity"`
	AveragePrice        *string  `json:"average_price"`
	QuotedCurrency      string   `json:"quoted_currency"`
	BaseCurrency        string   `json:"base_currency"`
	TickSize            string   `json:"tick_size"`
	Disabled            bool     `json:"disabled"`
	MarginEnabled       bool     `json:"margin_enabled"`
	CFDEnabled          bool     `json:"cfd_enabled"`
	PerpetualEnabled    bool     `json:"perpetual_enabled"`
	LastEventTimestamp  string   `json:"last_event_timestamp"`
	Timestamp           string   `json:"timestamp"`
	MultiplierUp        string   `json:"multiplier_up"`
	MultiplierDown      string   `json:"multiplier_down"`
	AverageTimeInterval *int     `json:"average_time_interval"`
}

// LiquidConnectionEstablished is the event name used when the connection is established.s
const LiquidConnectionEstablished = "pusher:connection_established"

// LiquidEventSubscribe is the event name used on subscription.
const LiquidEventSubscribe = "pusher:subscribe"

// LiquidEventSubscriptionSucceeded is the event name used when subscription is succeeded.
const LiquidEventSubscriptionSucceeded = "pusher_internal:subscription_succeeded"

// LiquidEventCreated is the event name used when a entity is newly created on a channel.
const LiquidEventCreated = "created"

// LiquidEventUpdated is the event name used when a entity is updated on a channel.
const LiquidEventUpdated = "updated"

// LiquidMessageRoot is root structure of WebSocket message from Liquid.
type LiquidMessageRoot struct {
	Channel *string         `json:"channel,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
	Event   string          `json:"event"`
}

// LiquidExecution is trade execution message from Liquid.
type LiquidExecution struct {
	ID        int64   `json:"id"`
	Quantity  float64 `json:"quantity"`
	Price     float64 `json:"price"`
	TakerSide string  `json:"taker_side"`
	CreatedAt int     `json:"created_at"`
}

// LiquidConnectionEstablishedData is data of connection established event.
type LiquidConnectionEstablishedData struct {
	ActivityTimeout int    `json:"activity_timeout"`
	SocketID        string `json:"socket_id"`
}

// LiquidSubscribeData is data of subscribe event.
type LiquidSubscribeData struct {
	Channel string `json:"channel"`
}
