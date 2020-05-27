package jsondef

// BitmexOrderBookL2 is auto-generated
type BitmexOrderBookL2 struct {
	Pair  string  `json:"pair"`
	Price float64 `json:"price"`
	ID    int64   `json:"id"`
	Size  float64 `json:"size"`
}

// TypeDefBitmexOrderBookL2 is auto-generated
var TypeDefBitmexOrderBookL2 = []byte("{\"pair\": \"pair\", \"price\": \"price\", \"id\": \"int\", \"size\": \"size\"}")

// BitmexTrade is auto-generated
type BitmexTrade struct {
	Pair            string   `json:"pair"`
	Price           float64  `json:"price"`
	Size            float64  `json:"size"`
	Timestamp       string   `json:"timestamp"`
	TrdMatchID      string   `json:"trdMatchId"`
	TickDirection   string   `json:"tickDirection"`
	GrossValue      *int64   `json:"grossValue"`
	HomeNotional    *float64 `json:"homeNotional"`
	ForeignNotional *float64 `json:"foreignNotional"`
}

// TypeDefBitmexTrade is auto-generated
var TypeDefBitmexTrade = []byte("{\"pair\": \"pair\", \"price\": \"price\", \"size\": \"size\", \"timestamp\": \"timestamp\", \"trdMatchId\": \"guid\", \"tickDirection\": \"string\", \"grossValue\": \"int\", \"homeNotional\": \"float\", \"foreignNotional\": \"float\"}")

// BitmexInstrument is auto-generated
type BitmexInstrument struct {
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

// TypeDefBitmexInstrument is auto-generated
var TypeDefBitmexInstrument = []byte("{\"symbol\": \"string\", \"rootSymbol\": \"string\", \"state\": \"string\", \"typ\": \"string\", \"listing\": \"timestamp\", \"front\": \"timestamp\", \"expiry\": \"timestamp\", \"settle\": \"timestamp\", \"relistInterval\": \"duration\", \"inverseLeg\": \"string\", \"sellLeg\": \"string\", \"buyLeg\": \"string\", \"optionStrikePcnt\": \"float\", \"optionStrikeRound\": \"float\", \"optionStrikePrice\": \"float\", \"optionMultiplier\": \"float\", \"positionCurrency\": \"string\", \"underlying\": \"string\", \"quoteCurrency\": \"string\", \"underlyingSymbol\": \"string\", \"reference\": \"string\", \"referenceSymbol\": \"string\", \"calcInterval\": \"duration\", \"publishInterval\": \"duration\", \"publishTime\": \"duration\", \"maxOrderQty\": \"int\", \"maxPrice\": \"float\", \"lotSize\": \"int\", \"tickSize\": \"float\", \"multiplier\": \"int\", \"settlCurrency\": \"string\", \"underlyingToPositionMultiplier\": \"int\", \"underlyingToSettleMultiplier\": \"int\", \"quoteToSettleMultiplier\": \"int\", \"isQuanto\": \"boolean\", \"isInverse\": \"boolean\", \"initMargin\": \"float\", \"maintMargin\": \"float\", \"riskLimit\": \"int\", \"riskStep\": \"int\", \"limit\": \"float\", \"capped\": \"boolean\", \"taxed\": \"boolean\", \"deleverage\": \"boolean\", \"makerFee\": \"float\", \"takerFee\": \"float\", \"settlementFee\": \"float\", \"insuranceFee\": \"float\", \"fundingBaseSymbol\": \"string\", \"fundingQuoteSymbol\": \"string\", \"fundingPremiumSymbol\": \"string\", \"fundingTimestamp\": \"timestamp\", \"fundingInterval\": \"duration\", \"fundingRate\": \"float\", \"indicativeFundingRate\": \"float\", \"rebalanceTimestamp\": \"timestamp\", \"rebalanceInterval\": \"duration\", \"openingTimestamp\": \"timestamp\", \"closingTimestamp\": \"timestamp\", \"sessionInterval\": \"duration\", \"prevClosePrice\": \"float\", \"limitDownPrice\": \"float\", \"limitUpPrice\": \"float\", \"bankruptLimitDownPrice\": \"float\", \"bankruptLimitUpPrice\": \"float\", \"prevTotalVolume\": \"int\", \"totalVolume\": \"int\", \"volume\": \"int\", \"volume24h\": \"int\", \"prevTotalTurnover\": \"int\", \"totalTurnover\": \"int\", \"turnover\": \"int\", \"turnover24h\": \"int\", \"homeNotional24h\": \"float\", \"foreignNotional24h\": \"float\", \"prevPrice24h\": \"float\", \"vwap\": \"float\", \"highPrice\": \"float\", \"lowPrice\": \"float\", \"lastPrice\": \"float\", \"lastPriceProtected\": \"float\", \"lastTickDirection\": \"string\", \"lastChangePcnt\": \"float\", \"bidPrice\": \"float\", \"midPrice\": \"float\", \"askPrice\": \"float\", \"impactBidPrice\": \"float\", \"impactMidPrice\": \"float\", \"impactAskPrice\": \"float\", \"hasLiquidity\": \"boolean\", \"openInterest\": \"int\", \"openValue\": \"int\", \"fairMethod\": \"string\", \"fairBasisRate\": \"float\", \"fairBasis\": \"float\", \"fairPrice\": \"float\", \"markMethod\": \"string\", \"markPrice\": \"float\", \"indicativeTaxRate\": \"float\", \"indicativeSettlePrice\": \"float\", \"optionUnderlyingPrice\": \"float\", \"settledPrice\": \"float\", \"timestamp\": \"timestamp\"}")

// BitmexInsurance is auto-generated
type BitmexInsurance struct {
	Currency      string `json:"currency"`
	Timestamp     string `json:"timestamp"`
	WalletBalance int64  `json:"walletBalance"`
}

// TypeDefBitmexInsurance is auto-generated
var TypeDefBitmexInsurance = []byte("{\"currency\": \"string\", \"timestamp\": \"timestamp\", \"walletBalance\": \"int\"}")

// BitmexFunding is auto-generated
type BitmexFunding struct {
	Timestamp        string  `json:"timestamp"`
	Symbol           string  `json:"symbol"`
	FundingInterval  string  `json:"fundingInterval"`
	FundingRate      float64 `json:"fundingRate"`
	FundingRateDaily float64 `json:"fundingRateDaily"`
}

// TypeDefBitmexFunding is auto-generated
var TypeDefBitmexFunding = []byte("{\"timestamp\": \"timestamp\", \"symbol\": \"string\", \"fundingInterval\": \"duration\", \"fundingRate\": \"float\", \"fundingRateDaily\": \"float\"}")

// BitmexSettlement is auto-generated
type BitmexSettlement struct {
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

// TypeDefBitmexSettlement is auto-generated
var TypeDefBitmexSettlement = []byte("{\"timestamp\": \"timestamp\", \"symbol\": \"string\", \"settlementType\": \"string\", \"settledPrice\": \"float\", \"optionStrikePrice\": \"float\", \"optionUnderlyingPrice\": \"float\", \"bankrupt\": \"int\", \"taxBase\": \"int\", \"taxRate\": \"float\"}")

// BitmexLiquidation is auto-generated
type BitmexLiquidation struct {
	OrderID   string  `json:"orderID"`
	Symbol    string  `json:"symbol"`
	Side      string  `json:"side"`
	Price     float64 `json:"price"`
	LeavesQty int64   `json:"leavesQty"`
}

// TypeDefBitmexLiquidation is auto-generated
var TypeDefBitmexLiquidation = []byte("{\"orderID\": \"guid\", \"symbol\": \"string\", \"side\": \"string\", \"price\": \"float\", \"leavesQty\": \"int\"}")
