package jsonstructs

import "encoding/json"

// BitmexFundingDataElement is a element of data of funding channel
type BitmexFundingDataElement struct {
	Timestamp        string  `json:"timestamp"`
	Symbol           string  `json:"symbol"`
	FundingInterval  string  `json:"fundingInterval"`
	FundingRate      float64 `json:"fundingRate"`
	FundingRateDaily float64 `json:"fundingRateDaily"`
}

// BitmexInsuranceDataElement is a element of data of insurance channel
type BitmexInsuranceDataElement struct {
	Currency      string `json:"currency"`
	Timestamp     string `json:"timestamp"`
	WalletBalance int64  `json:"walletBalance"`
}

// BitmexSettlementDataElement is a element of data of settlement channel from bitmex exchange
type BitmexSettlementDataElement struct {
	Timestamp             string  `json:"timestamp"`
	Symbol                string  `json:"symbol"`
	SettlementType        string  `json:"settlementType"`
	SettledPrice          float64 `json:"settledPrice"`
	OptionStrikePrice     float64 `json:"optionStrikePrice"`
	OptionUnderlyingPrice float64 `json:"optionUnderlyingPrice"`
	Bankrupt              int64   `json:"bankrupt"`
	TaxBase               int64   `json:"taxBase"`
	TaxRate               float64 `json:"taxRate"`
}

// BitmexLiquidationDataElement is a element of data of liquidation channel from bitmex exchange
type BitmexLiquidationDataElement struct {
	OrderID   string  `json:"orderId"`
	Symbol    string  `json:"symbol"`
	Side      string  `json:"side"`
	Price     float64 `json:"price"`
	LeavesQty int64   `json:"leavesQty"`
}

// BitmexInstrumentDataElem is a element of data of instrument channel from bitmex exchange
type BitmexInstrumentDataElem struct {
	Symbol                         string   `json:"symbol"`
	RootSymbol                     *string  `json:"rootSymbol"`
	State                          *string  `json:"state"`
	Typ                            *string  `json:"typ"`
	Listing                        *string  `json:"listing"`
	Front                          *string  `json:"front"`
	Expiry                         *string  `json:"expiry"`
	Settle                         *string  `json:"settle"`
	RelistInterval                 *string  `json:"relistInterval"`
	InverseLeg                     *string  `json:"inverseLeg"`
	SellLeg                        *string  `json:"sellLeg"`
	BuyLeg                         *string  `json:"buyLeg"`
	OptionStrikePcnt               *float64 `json:"optionStrikePcnt"`
	OptionStrikeRound              *float64 `json:"optionStrikeRound"`
	OptionStrikePrice              *float64 `json:"optionStrikePrice"`
	OptionMultiplier               *float64 `json:"optionMultiplier"`
	PositionCurrency               *string  `json:"positionCurrency"`
	Underlying                     *string  `json:"underlying"`
	QuoteCurrency                  *string  `json:"quoteCurrency"`
	UnderlyingSymbol               *string  `json:"underlyingSymbol"`
	Reference                      *string  `json:"reference"`
	ReferenceSymbol                *string  `json:"referenceSymbol"`
	CalcInterval                   *string  `json:"calcInterval"`
	PublishInterval                *string  `json:"publishInterval"`
	PublishTime                    *string  `json:"publishTime"`
	MaxOrderQty                    *int64   `json:"maxOrderQty"`
	MaxPrice                       *float64 `json:"maxPrice"`
	LotSize                        *int64   `json:"lotSize"`
	TickSize                       *float64 `json:"tickSize"`
	Multiplier                     *int64   `json:"multiplier"`
	SettlCurrency                  *string  `json:"settlCurrency"`
	UnderlyingToPositionMultiplier *int64   `json:"underlyingToPositionMultiplier"`
	UnderlyingToSettleMultiplier   *int64   `json:"underlyingToSettleMultiplier"`
	QuoteToSettleMultiplier        *int64   `json:"quoteToSettleMultiplier"`
	IsQuanto                       *bool    `json:"isQuanto"`
	IsInverse                      *bool    `json:"isInverse"`
	InitMargin                     *float64 `json:"initMargin"`
	MaintMargin                    *float64 `json:"maintMargin"`
	RiskLimit                      *int64   `json:"riskLimit"`
	RiskStep                       *int64   `json:"riskStep"`
	Limit                          *float64 `json:"limit"`
	Capped                         *bool    `json:"capped"`
	Taxed                          *bool    `json:"taxed"`
	Deleverage                     *bool    `json:"deleverage"`
	MakerFee                       *float64 `json:"makerFee"`
	TakerFee                       *float64 `json:"takerFee"`
	SettlementFee                  *float64 `json:"settlementFee"`
	InsuranceFee                   *float64 `json:"insuranceFee"`
	FundingBaseSymbol              *string  `json:"fundingBaseSymbol"`
	FundingQuoteSymbol             *string  `json:"fundingQuoteSymbol"`
	FundingPremiumSymbol           *string  `json:"fundingPremiumSymbol"`
	FundingTimestamp               *string  `json:"fundingTimestamp"`
	FundingInterval                *string  `json:"fundingInterval"`
	FundingRate                    *float64 `json:"fundingRate"`
	IndicativeFundingRate          *float64 `json:"indicativeFundingRate"`
	RebalanceTimestamp             *string  `json:"rebalanceTimestamp"`
	RebalanceInterval              *string  `json:"rebalanceInterval"`
	OpeningTimestamp               *string  `json:"openingTimestamp"`
	ClosingTimestamp               *string  `json:"closingTimestamp"`
	SessionInterval                *string  `json:"sessionInterval"`
	PrevClosePrice                 *float64 `json:"prevClosePrice"`
	LimitDownPrice                 *float64 `json:"limitDownPrice"`
	LimitUpPrice                   *float64 `json:"limitUpPrice"`
	BankruptLimitDownPrice         *float64 `json:"bankruptLimitDownPrice"`
	BankruptLimitUpPrice           *float64 `json:"bankruptLimitUpPrice"`
	PrevTotalVolume                *int64   `json:"prevTotalVolume"`
	TotalVolume                    *int64   `json:"totalVolume"`
	Volume                         *int64   `json:"volume"`
	Volume24h                      *int64   `json:"volume24h"`
	PrevTotalTurnover              *int64   `json:"prevTotalTurnover"`
	TotalTurnover                  *int64   `json:"totalTurnover"`
	Turnover                       *int64   `json:"turnover"`
	Turnover24h                    *int64   `json:"turnover24h"`
	HomeNotional24h                *float64 `json:"homeNotional24h"`
	ForeignNotional24h             *float64 `json:"foreignNotional24h"`
	PrevPrice24h                   *float64 `json:"prevPrice24h"`
	Vwap                           *float64 `json:"vwap"`
	HighPrice                      *float64 `json:"highPrice"`
	LowPrice                       *float64 `json:"lowPrice"`
	LastPrice                      *float64 `json:"lastPrice"`
	LastPriceProtected             *float64 `json:"lastPriceProtected"`
	LastTickDirection              *string  `json:"lastTickDirection"`
	LastChangePcnt                 *float64 `json:"lastChangePcnt"`
	BidPrice                       *float64 `json:"bidPrice"`
	MidPrice                       *float64 `json:"midPrice"`
	AskPrice                       *float64 `json:"askPrice"`
	ImpactBidPrice                 *float64 `json:"impactBidPrice"`
	ImpactMidPrice                 *float64 `json:"impactMidPrice"`
	ImpactAskPrice                 *float64 `json:"impactAskPrice"`
	HasLiquidity                   *bool    `json:"hasLiquidity"`
	OpenInterest                   *int64   `json:"openInterest"`
	OpenValue                      *int64   `json:"openValue"`
	FairMethod                     *string  `json:"fairMethod"`
	FairBasisRate                  *float64 `json:"fairBasisRate"`
	FairBasis                      *float64 `json:"fairBasis"`
	FairPrice                      *float64 `json:"fairPrice"`
	MarkMethod                     *string  `json:"markMethod"`
	MarkPrice                      *float64 `json:"markPrice"`
	IndicativeTaxRate              *float64 `json:"indicativeTaxRate"`
	IndicativeSettlePrice          *float64 `json:"indicativeSettlePrice"`
	OptionUnderlyingPrice          *float64 `json:"optionUnderlyingPrice"`
	SettledPrice                   *float64 `json:"settledPrice"`
	Timestamp                      string   `json:"timestamp"`
}

// BitmexOrderBookL2DataElement is orderbook element
type BitmexOrderBookL2DataElement struct {
	Symbol string  `json:"symbol"`
	ID     int64   `json:"id"`
	Side   string  `json:"side"`
	Price  float64 `json:"price"`
	Size   uint64  `json:"size"`
}

// BitmexTradeDataElement individual trade order
type BitmexTradeDataElement struct {
	Timestamp     string  `json:"timestamp"`
	Symbol        string  `json:"symbol"`
	Side          string  `json:"side"`
	Size          uint64  `json:"size"`
	Price         float64 `json:"price"`
	TickDirection string  `json:"tickDirection"`
	TradeMatchID  string  `json:"trdMatchID"`
	// those 3 could be null
	GrossValue      *int64   `json:"grossValue"`
	HomeNotional    *float64 `json:"homeNotional"`
	ForeignNotional *float64 `json:"foreignNotional"`
}

// BitmexRoot is the root structure of bitmex exchange message
type BitmexRoot struct {
	Table  string          `json:"table"`
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
	// present only if this is a infomation message
	Info *string `json:"info"`
	// present only if this is a error message (error to subscribe request)
	Error *string `json:"error"`
}

// BitmexSubscribe is struct for subscribe message from bitmex
type BitmexSubscribe struct {
	Success   bool   `json:"success"`
	Subscribe string `json:"subscribe"`
}

// Initialize BitmexSubscribe
func (s *BitmexSubscribe) Initialize() {
	s.Success = true
}

// BitmexStateSubscribed is a list of subscribed channels listed in state line in dataset
type BitmexStateSubscribed []string
