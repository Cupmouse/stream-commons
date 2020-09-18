package formatter

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/exchangedataset/streamcommons"

	"github.com/exchangedataset/streamcommons/formatter/jsondef"
	"github.com/exchangedataset/streamcommons/jsonstructs"
)

var bitmexDurationBaseTime time.Time

func bitmexParseTimestamp(timestamp *string) (*string, error) {
	if timestamp != nil {
		timestampTime, serr := time.Parse(time.RFC3339Nano, *timestamp)
		if serr != nil {
			return nil, fmt.Errorf("bitmexParseTimestamp: %v", serr)
		}
		result := strconv.FormatInt(timestampTime.UnixNano(), 10)
		return &result, nil
	}
	return nil, nil
}

func bitmexParseDuration(duration *string) (*string, error) {
	if duration != nil {
		durationTime, serr := time.Parse(time.RFC3339Nano, *duration)
		if serr != nil {
			return nil, fmt.Errorf("bitmexParseDuration: %v", serr)
		}
		result := strconv.FormatInt(durationTime.Sub(bitmexDurationBaseTime).Nanoseconds(), 10)
		return &result, nil
	}
	return nil, nil
}

// bitmexFormatter formats message from bitmex
type bitmexFormatter struct {
	// Target channels are used to support symbol-wise filtering and typedef formatting a for start line.
	targets map[string]bool
}

// FormatStart returns empty slice.
func (f *bitmexFormatter) FormatStart(urlStr string) ([]Result, error) {
	return make([]Result, 0), nil
}

func (f *bitmexFormatter) formatOrderBookL2(dataRaw json.RawMessage) (ret []Result, err error) {
	orders := make([]jsonstructs.BitmexOrderBookL2DataElement, 0, 10)
	serr := json.Unmarshal(dataRaw, &orders)
	if serr != nil {
		err = fmt.Errorf("formatOrderBookL2: BitmexOrderBookL2DataElement: %v", serr)
		return
	}

	ret = make([]Result, len(orders))
	for i, order := range orders {
		size := float64(order.Size)
		marshaled, serr := json.Marshal(jsondef.BitmexOrderBookL2{
			Symbol: order.Symbol,
			Price:  order.Price,
			ID:     order.ID,
			// Our side string format are the same as Bitmex
			Side: order.Side,
			Size: size,
		})
		if serr != nil {
			err = fmt.Errorf("formatOrderBookL2: BitmexOrderBookL2: %v", serr)
			return
		}
		ret[i] = Result{
			Channel: streamcommons.BitmexChannelOrderBookL2 + "_" + order.Symbol,
			Message: marshaled,
		}
	}
	return
}

func (f *bitmexFormatter) formatTrade(dataRaw json.RawMessage) (ret []Result, err error) {
	orders := make([]jsonstructs.BitmexTradeDataElement, 0, 10)
	serr := json.Unmarshal(dataRaw, &orders)
	if serr != nil {
		err = fmt.Errorf("formatTrade: BitmexTradeDataElement: %v", serr)
		return
	}
	ret = make([]Result, len(orders))
	for i, elem := range orders {
		size := float64(elem.Size)
		timestamp, serr := bitmexParseTimestamp(&elem.Timestamp)
		marshaled, serr := json.Marshal(jsondef.BitmexTrade{
			Symbol:          elem.Symbol,
			Price:           elem.Price,
			Side:            elem.Side,
			Size:            size,
			Timestamp:       *timestamp,
			TrdMatchID:      elem.TradeMatchID,
			TickDirection:   elem.TickDirection,
			GrossValue:      elem.GrossValue,
			HomeNotional:    elem.HomeNotional,
			ForeignNotional: elem.ForeignNotional,
		})
		if serr != nil {
			err = fmt.Errorf("formatTrade: BitmexTrade: %v", serr)
			return
		}
		ret[i] = Result{
			Channel: streamcommons.BitmexChannelTrade + "_" + elem.Symbol,
			Message: marshaled,
		}
	}
	return
}

func (f *bitmexFormatter) formatInstrument(dataRaw json.RawMessage) (ret []Result, err error) {
	instruments := make([]jsonstructs.BitmexInstrumentDataElem, 0, 10)
	serr := json.Unmarshal(dataRaw, &instruments)
	if serr != nil {
		err = fmt.Errorf("formatInstrument: dataRaw: %v", serr)
		return
	}

	ret = make([]Result, len(instruments))
	for i, elem := range instruments {
		relistInterval, serr := bitmexParseDuration(elem.RelistInterval)
		if serr != nil {
			err = fmt.Errorf("formatInstrument: relistInterval: %v", serr)
			return
		}
		calcInterval, serr := bitmexParseDuration(elem.CalcInterval)
		if serr != nil {
			err = fmt.Errorf("formatInstrument: calcInterval: %v", serr)
			return
		}
		publishInterval, serr := bitmexParseDuration(elem.PublishInterval)
		if serr != nil {
			err = fmt.Errorf("formatInstrument: publishInterval: %v", serr)
			return
		}
		publishTime, serr := bitmexParseDuration(elem.PublishTime)
		if serr != nil {
			err = fmt.Errorf("formatInstrument: publishTime: %v", serr)
			return
		}
		fundingInterval, serr := bitmexParseDuration(elem.FundingInterval)
		if serr != nil {
			err = fmt.Errorf("formatInstrument: fundingInterval: %v", serr)
			return
		}
		rebalanceInterval, serr := bitmexParseDuration(elem.RebalanceInterval)
		if serr != nil {
			err = fmt.Errorf("formatInstrument: rebalanceInterval: %v", serr)
			return
		}
		sessionInterval, serr := bitmexParseDuration(elem.SessionInterval)
		if serr != nil {
			err = fmt.Errorf("formatInstrument: sessionInterval: %v", serr)
			return
		}
		listing, serr := bitmexParseTimestamp(elem.Listing)
		if serr != nil {
			err = fmt.Errorf("formatInstrument: listing: %v", serr)
			return
		}
		front, serr := bitmexParseTimestamp(elem.Front)
		if serr != nil {
			err = fmt.Errorf("formatInstrument: front: %v", serr)
			return
		}
		expiry, serr := bitmexParseTimestamp(elem.Expiry)
		if serr != nil {
			err = fmt.Errorf("formatInstrument: expiry: %v", serr)
			return
		}
		settle, serr := bitmexParseTimestamp(elem.Settle)
		if serr != nil {
			err = fmt.Errorf("formatInstrument: settle: %v", serr)
			return
		}
		closingTimestamp, serr := bitmexParseTimestamp(elem.ClosingTimestamp)
		if serr != nil {
			err = fmt.Errorf("formatInstrument: closingTimestamp: %v", serr)
			return
		}
		fundingTimestamp, serr := bitmexParseTimestamp(elem.FundingTimestamp)
		if serr != nil {
			err = fmt.Errorf("formatInstrument: fundingTimestamp: %v", serr)
			return
		}
		openingTimestamp, serr := bitmexParseTimestamp(elem.OpeningTimestamp)
		if serr != nil {
			err = fmt.Errorf("formatInstrument: openingTimestamp: %v", serr)
			return
		}
		rebalanceTimestamp, serr := bitmexParseTimestamp(elem.RebalanceTimestamp)
		if serr != nil {
			err = fmt.Errorf("formatInstrument: rebalanceTimestamp: %v", serr)
			return
		}
		timestamp, serr := bitmexParseTimestamp(&elem.Timestamp)
		if serr != nil {
			err = fmt.Errorf("formatInstrument: timestamp: %v", serr)
			return
		}
		marshaled, serr := json.Marshal(jsondef.BitmexInstrument{
			Symbol:                         elem.Symbol,
			RootSymbol:                     elem.RootSymbol,
			State:                          elem.State,
			Typ:                            elem.Typ,
			Listing:                        listing,
			Front:                          front,
			Expiry:                         expiry,
			Settle:                         settle,
			RelistInterval:                 relistInterval,
			InverseLeg:                     elem.InverseLeg,
			SellLeg:                        elem.SellLeg,
			BuyLeg:                         elem.BuyLeg,
			OptionStrikePcnt:               elem.OptionStrikePcnt,
			OptionStrikeRound:              elem.OptionStrikeRound,
			OptionStrikePrice:              elem.OptionStrikePrice,
			OptionMultiplier:               elem.OptionMultiplier,
			PositionCurrency:               elem.PositionCurrency,
			Underlying:                     elem.Underlying,
			QuoteCurrency:                  elem.QuoteCurrency,
			UnderlyingSymbol:               elem.UnderlyingSymbol,
			Reference:                      elem.Reference,
			ReferenceSymbol:                elem.ReferenceSymbol,
			CalcInterval:                   calcInterval,
			PublishInterval:                publishInterval,
			PublishTime:                    publishTime,
			MaxOrderQty:                    elem.MaxOrderQty,
			MaxPrice:                       elem.MaxPrice,
			LotSize:                        elem.LotSize,
			TickSize:                       elem.TickSize,
			Multiplier:                     elem.Multiplier,
			SettlCurrency:                  elem.SettlCurrency,
			UnderlyingToPositionMultiplier: elem.UnderlyingToPositionMultiplier,
			UnderlyingToSettleMultiplier:   elem.UnderlyingToSettleMultiplier,
			QuoteToSettleMultiplier:        elem.QuoteToSettleMultiplier,
			IsQuanto:                       elem.IsQuanto,
			IsInverse:                      elem.IsInverse,
			InitMargin:                     elem.InitMargin,
			MaintMargin:                    elem.MaintMargin,
			RiskLimit:                      elem.RiskLimit,
			RiskStep:                       elem.RiskStep,
			Limit:                          elem.Limit,
			Capped:                         elem.Capped,
			Taxed:                          elem.Taxed,
			Deleverage:                     elem.Deleverage,
			MakerFee:                       elem.MakerFee,
			TakerFee:                       elem.TakerFee,
			SettlementFee:                  elem.SettlementFee,
			InsuranceFee:                   elem.InsuranceFee,
			FundingBaseSymbol:              elem.FundingBaseSymbol,
			FundingQuoteSymbol:             elem.FundingQuoteSymbol,
			FundingPremiumSymbol:           elem.FundingPremiumSymbol,
			FundingTimestamp:               fundingTimestamp,
			FundingInterval:                fundingInterval,
			FundingRate:                    elem.FundingRate,
			IndicativeFundingRate:          elem.IndicativeFundingRate,
			RebalanceTimestamp:             rebalanceTimestamp,
			RebalanceInterval:              rebalanceInterval,
			OpeningTimestamp:               openingTimestamp,
			ClosingTimestamp:               closingTimestamp,
			SessionInterval:                sessionInterval,
			PrevClosePrice:                 elem.PrevClosePrice,
			LimitDownPrice:                 elem.LimitDownPrice,
			LimitUpPrice:                   elem.LimitUpPrice,
			BankruptLimitDownPrice:         elem.BankruptLimitDownPrice,
			BankruptLimitUpPrice:           elem.BankruptLimitUpPrice,
			PrevTotalVolume:                elem.PrevTotalVolume,
			TotalVolume:                    elem.TotalVolume,
			Volume:                         elem.Volume,
			Volume24h:                      elem.Volume24h,
			PrevTotalTurnover:              elem.PrevTotalTurnover,
			TotalTurnover:                  elem.TotalTurnover,
			Turnover:                       elem.Turnover,
			Turnover24h:                    elem.Turnover24h,
			HomeNotional24h:                elem.HomeNotional24h,
			ForeignNotional24h:             elem.ForeignNotional24h,
			PrevPrice24h:                   elem.PrevPrice24h,
			Vwap:                           elem.Vwap,
			HighPrice:                      elem.HighPrice,
			LowPrice:                       elem.LowPrice,
			LastPrice:                      elem.LastPrice,
			LastPriceProtected:             elem.LastPriceProtected,
			LastTickDirection:              elem.LastTickDirection,
			LastChangePcnt:                 elem.LastChangePcnt,
			BidPrice:                       elem.BidPrice,
			MidPrice:                       elem.MidPrice,
			AskPrice:                       elem.AskPrice,
			ImpactBidPrice:                 elem.ImpactBidPrice,
			ImpactMidPrice:                 elem.ImpactMidPrice,
			ImpactAskPrice:                 elem.ImpactAskPrice,
			HasLiquidity:                   elem.HasLiquidity,
			OpenInterest:                   elem.OpenInterest,
			OpenValue:                      elem.OpenValue,
			FairMethod:                     elem.FairMethod,
			FairBasisRate:                  elem.FairBasisRate,
			FairBasis:                      elem.FairBasis,
			FairPrice:                      elem.FairPrice,
			MarkMethod:                     elem.MarkMethod,
			MarkPrice:                      elem.MarkPrice,
			IndicativeTaxRate:              elem.IndicativeTaxRate,
			IndicativeSettlePrice:          elem.IndicativeSettlePrice,
			OptionUnderlyingPrice:          elem.OptionUnderlyingPrice,
			SettledPrice:                   elem.SettledPrice,
			Timestamp:                      *timestamp,
		})
		if serr != nil {
			err = fmt.Errorf("formatInstrument: BitmexInstrument: %v", serr)
			return
		}

		ret[i] = Result{
			Channel: streamcommons.BitmexChannelInstrument + "_" + elem.Symbol,
			Message: marshaled,
		}
	}
	return
}

func (f *bitmexFormatter) formatLiquidation(dataRaw json.RawMessage) (ret []Result, err error) {
	liquidations := make([]jsonstructs.BitmexLiquidationDataElement, 0, 10)
	serr := json.Unmarshal(dataRaw, &liquidations)
	if serr != nil {
		err = fmt.Errorf("formatLiquidation: dataRaw: %v", serr)
		return
	}
	ret = make([]Result, len(liquidations))
	for i, elem := range liquidations {
		// Note: this structure is completely the same as what bitmex sends us...
		marshaled, serr := json.Marshal(jsondef.BitmexLiquidation{
			OrderID:   elem.OrderID,
			Symbol:    elem.Symbol,
			Side:      elem.Side,
			Price:     elem.Price,
			LeavesQty: elem.LeavesQty,
		})
		if serr != nil {
			err = fmt.Errorf("formatLiquidation: BitmexLiquidation: %v", serr)
			return
		}
		ret[i] = Result{
			Channel: streamcommons.BitmexChannelLiquidation + "_" + elem.Symbol,
			Message: marshaled,
		}
	}
	return
}

func (f *bitmexFormatter) formatSettlement(dataRaw json.RawMessage) (ret []Result, err error) {
	settlements := make([]jsonstructs.BitmexSettlementDataElement, 0, 10)
	serr := json.Unmarshal(dataRaw, &settlements)
	if serr != nil {
		err = fmt.Errorf("formatSettlement: dataRaw: %v", serr)
		return
	}
	ret = make([]Result, len(settlements))
	for i, elem := range settlements {
		timestamp, serr := bitmexParseTimestamp(&elem.Timestamp)
		marshaled, serr := json.Marshal(jsondef.BitmexSettlement{
			Timestamp:             *timestamp,
			Symbol:                elem.Symbol,
			SettlementType:        elem.SettlementType,
			SettledPrice:          elem.SettledPrice,
			OptionStrikePrice:     elem.OptionStrikePrice,
			OptionUnderlyingPrice: elem.OptionUnderlyingPrice,
			Bankrupt:              elem.Bankrupt,
			TaxBase:               elem.TaxBase,
			TaxRate:               elem.TaxRate,
		})
		if serr != nil {
			err = fmt.Errorf("formatSettlement: BitmexSettlement: %v", serr)
			return
		}
		ret[i] = Result{
			Channel: streamcommons.BitmexChannelSettlement + "_" + elem.Symbol,
			Message: marshaled,
		}
	}
	return
}

func (f *bitmexFormatter) formatInsurance(dataRaw json.RawMessage) (ret []Result, err error) {
	insurances := make([]jsonstructs.BitmexInsuranceDataElement, 0, 10)
	serr := json.Unmarshal(dataRaw, &insurances)
	if serr != nil {
		err = fmt.Errorf("formatInsurance: BitmexInsuranceDataElement: %v", serr)
		return
	}
	ret = make([]Result, len(insurances))
	for i, elem := range insurances {
		timestamp, serr := bitmexParseTimestamp(&elem.Timestamp)
		marshaled, serr := json.Marshal(jsondef.BitmexInsurance{
			Currency:      elem.Currency,
			Timestamp:     *timestamp,
			WalletBalance: elem.WalletBalance,
		})
		if serr != nil {
			err = fmt.Errorf("formatInsurance: BitmexInsurance: %v", serr)
			return
		}
		ret[i] = Result{
			Channel: streamcommons.BitmexChannelInsurance + "_" + elem.Currency,
			Message: marshaled,
		}
	}
	return
}

func (f *bitmexFormatter) formatFunding(dataRaw json.RawMessage) (ret []Result, err error) {
	fundings := make([]jsonstructs.BitmexFundingDataElement, 0, 10)
	serr := json.Unmarshal(dataRaw, &fundings)
	if serr != nil {
		err = fmt.Errorf("formatFunding: BitmexFundingDataElement: %v", serr)
		return
	}
	ret = make([]Result, len(fundings))
	for i, elem := range fundings {
		timestamp, serr := bitmexParseTimestamp(&elem.Timestamp)
		if serr != nil {
			err = fmt.Errorf("formatFunding: timestamp: %v", serr)
		}
		fundingInterval, serr := bitmexParseDuration(&elem.FundingInterval)
		if serr != nil {
			err = fmt.Errorf("formatFunding: fundingInterval: %v", serr)
		}
		marshaled, serr := json.Marshal(jsondef.BitmexFunding{
			Timestamp:        *timestamp,
			Symbol:           elem.Symbol,
			FundingInterval:  *fundingInterval,
			FundingRate:      elem.FundingRate,
			FundingRateDaily: elem.FundingRateDaily,
		})
		if serr != nil {
			err = fmt.Errorf("formatFunding: BitmexFunding: %v", serr)
			return
		}
		ret[i] = Result{
			Channel: streamcommons.BitmexChannelFunding + "_" + elem.Symbol,
			Message: marshaled,
		}
	}
	return
}

// FormatMessage formats incoming message given and returns formatted strings
func (f *bitmexFormatter) FormatMessage(channel string, line []byte) (ret []Result, err error) {
	subscribed := jsonstructs.BitmexSubscribe{}
	serr := json.Unmarshal(line, &subscribed)
	if serr != nil {
		err = fmt.Errorf("FormatMessage: BitmexSubscribe: %v", serr)
		return
	}
	if subscribed.Success {
		// This is a response to subscription
		// Return header row
		var typedef []byte
		switch channel {
		case streamcommons.BitmexChannelOrderBookL2:
			typedef = jsondef.TypeDefBitmexOrderBookL2
		case streamcommons.BitmexChannelTrade:
			typedef = jsondef.TypeDefBitmexTrade
		case streamcommons.BitmexChannelInstrument:
			typedef = jsondef.TypeDefBitmexInstrument
		case streamcommons.BitmexChannelLiquidation:
			typedef = jsondef.TypeDefBitmexLiquidation
		case streamcommons.BitmexChannelSettlement:
			typedef = jsondef.TypeDefBitmexSettlement
		case streamcommons.BitmexChannelInsurance:
			typedef = jsondef.TypeDefBitmexInsurance
		case streamcommons.BitmexChannelFunding:
			typedef = jsondef.TypeDefBitmexFunding
		default:
			err = fmt.Errorf("FormatMessage: json not supported: %s", channel)
			return
		}
		ret = make([]Result, 0, len(f.targets))
		for target := range f.targets {
			// Is this target a symbol-wise of raw channel?
			if strings.HasPrefix(target, channel) {
				// Duplicate typedef to symbol-wise channel
				ret = append(ret, Result{
					Channel: target,
					Message: typedef,
				})
			}
		}
		return
	}

	root := new(jsonstructs.BitmexRoot)
	serr = json.Unmarshal(line, root)
	if serr != nil {
		err = fmt.Errorf("FormatMessage: BitmexRoot: %v", serr)
		return
	}

	var sret []Result
	switch channel {
	case streamcommons.BitmexChannelOrderBookL2:
		sret, err = f.formatOrderBookL2(root.Data)
	case streamcommons.BitmexChannelTrade:
		sret, err = f.formatTrade(root.Data)
	case streamcommons.BitmexChannelInstrument:
		sret, err = f.formatInstrument(root.Data)
	case streamcommons.BitmexChannelLiquidation:
		sret, err = f.formatLiquidation(root.Data)
	case streamcommons.BitmexChannelSettlement:
		sret, err = f.formatSettlement(root.Data)
	case streamcommons.BitmexChannelInsurance:
		sret, err = f.formatInsurance(root.Data)
	case streamcommons.BitmexChannelFunding:
		sret, err = f.formatFunding(root.Data)
	default:
		err = fmt.Errorf("FormatMessage: json unsupported: %s", channel)
		return
	}
	// Apply post filter
	ret = make([]Result, 0, len(sret))
	for _, r := range sret {
		if _, ok := f.targets[r.Channel]; ok {
			ret = append(ret, r)
		}
	}
	return
}

// IsSupported returns true if given channel is supported to be formatted using this formatter
func (f *bitmexFormatter) IsSupported(channel string) bool {
	if pos := strings.IndexRune(channel, '_'); pos != -1 {
		// Symbol-wise channels are also supported
		channel = channel[:pos]
	}
	return channel == "orderBookL2" || channel == "trade" || channel == "instrument" ||
		channel == "liquidation" || channel == "settlement" || channel == "insurance" ||
		channel == "funding"
}

func init() {
	var serr error
	bitmexDurationBaseTime, serr = time.Parse(time.RFC3339Nano, "2000-01-01T00:00:00.000Z")
	if serr != nil {
		panic(fmt.Sprintf("init durationBaseTime: %v", serr))
	}
}

func newBitmexFormatter(symbolWiseTargets []string) *bitmexFormatter {
	f := new(bitmexFormatter)
	f.targets = make(map[string]bool)
	for _, ch := range symbolWiseTargets {
		f.targets[ch] = true
	}
	return f
}
