package formatter

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/exchangedataset/streamcommons"
	"github.com/exchangedataset/streamcommons/formatter/jsondef"
	"github.com/exchangedataset/streamcommons/jsonstructs"
)

// binanceFormatter is json formatter for Binance.
type binanceFormatter struct{}

// FormatStart formats start line (URL) and returns the array of known subscribed channel in case the server won't
// tell the client what channels are successfully subscribed.
func (f *binanceFormatter) FormatStart(urlStr string) (formatted []Result, err error) {
	u, serr := url.Parse(string(urlStr))
	if serr != nil {
		return nil, fmt.Errorf("FormatStart: %v", serr)
	}
	q := u.Query()
	streams := q.Get("streams")
	channels := strings.Split(streams, "/")
	formatted = make([]Result, len(channels))
	for i, ch := range channels {
		_, stream, serr := streamcommons.BinanceDecomposeChannel(ch)
		if serr != nil {
			err = fmt.Errorf("FormatStart: %v", serr)
			return
		}
		switch stream {
		case streamcommons.BinanceStreamDepth:
			formatted[i] = Result{
				Channel: ch,
				Message: jsondef.TypeDefBinanceDepth,
			}
		case streamcommons.BinanceStreamTrade:
			formatted[i] = Result{
				Channel: ch,
				Message: jsondef.TypeDefBinanceTrade,
			}
		case streamcommons.BinanceStreamTicker:
			formatted[i] = Result{
				Channel: ch,
				Message: jsondef.TypeDefBinanceTicker,
			}
		case streamcommons.BinanceStreamRESTDepth:
			formatted[i] = Result{
				Channel: ch,
				Message: jsondef.TypeDefBinanceRestDepth,
			}
		default:
			err = fmt.Errorf("FormatStart: channel not supported: %s", ch)
			return
		}
	}
	return formatted, nil
}

// FormatMessage formats messages from server.
func (f *binanceFormatter) FormatMessage(channel string, line []byte) (formatted []Result, err error) {
	symbol, stream, serr := streamcommons.BinanceDecomposeChannel(channel)
	if serr != nil {
		err = fmt.Errorf("FormatMessage: %v", serr)
		return
	}
	subscribed := new(jsonstructs.BinanceSubscribeResponse)
	serr = json.Unmarshal(line, subscribed)
	if serr != nil {
		err = fmt.Errorf("FormatMessage: line: %v", serr)
		return
	}
	if subscribed.ID != 0 {
		// Subscribe message
		formatted = make([]Result, 1)
		switch stream {
		case streamcommons.BinanceStreamDepth:
			formatted[0] = Result{
				Channel: channel,
				Message: jsondef.TypeDefBinanceDepth,
			}
		case streamcommons.BinanceStreamTrade:
			formatted[0] = Result{
				Channel: channel,
				Message: jsondef.TypeDefBinanceTrade,
			}
		case streamcommons.BinanceStreamRESTDepth:
			formatted[0] = Result{
				Channel: channel,
				Message: jsondef.TypeDefBinanceRestDepth,
			}
		case streamcommons.BinanceStreamTicker:
			formatted[0] = Result{
				Channel: channel,
				Message: jsondef.TypeDefBinanceTicker,
			}
		default:
			err = fmt.Errorf("FormatMessage: channel not supported: %s", channel)
		}
		return
	}
	switch stream {
	case streamcommons.BinanceStreamDepth:
		root := new(jsonstructs.BinanceReponseRoot)
		serr := json.Unmarshal(line, root)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: depth BinanceReponseRoot: %v", serr)
			return
		}
		depth := new(jsonstructs.BinanceDepthStream)
		serr = json.Unmarshal(root.Data, depth)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: depth BinanceDepthStream: %v", serr)
			return
		}
		eventTime := strconv.FormatInt(depth.EventTime*int64(time.Millisecond), 10)
		formatted = make([]Result, len(depth.Asks)+len(depth.Bids))
		i := 0
		for _, order := range depth.Asks {
			fo := new(jsondef.BinanceDepth)
			fo.Symbol = depth.Symbol
			fo.EventTime = eventTime
			fo.Price, serr = strconv.ParseFloat(order[0], 64)
			if serr != nil {
				err = fmt.Errorf("FormatMessage: depth ask price: %v", serr)
				return
			}
			size, serr := strconv.ParseFloat(order[1], 64)
			if serr != nil {
				err = fmt.Errorf("FormatMessage: depth ask size: %v", serr)
				return
			}
			// Negative size for sell order
			fo.Size = -size
			mfo, serr := json.Marshal(fo)
			if serr != nil {
				err = fmt.Errorf("FormatMessage: depth ask BinanceDepth: %v", serr)
				return
			}
			formatted[i] = Result{
				Channel: channel,
				Message: mfo,
			}
			i++
		}
		for _, order := range depth.Bids {
			fo := new(jsondef.BinanceDepth)
			fo.Symbol = depth.Symbol
			fo.EventTime = eventTime
			fo.Price, serr = strconv.ParseFloat(order[0], 64)
			if serr != nil {
				err = fmt.Errorf("FormatMessage: depth bid price: %v", serr)
				return
			}
			fo.Size, serr = strconv.ParseFloat(order[1], 64)
			if serr != nil {
				err = fmt.Errorf("FormatMessage: depth bid size: %v", serr)
				return
			}
			mfo, serr := json.Marshal(fo)
			if serr != nil {
				err = fmt.Errorf("FormatMessage: depth bid BinanceDepth: %v", serr)
				return
			}
			formatted[i] = Result{
				Channel: channel,
				Message: mfo,
			}
			i++
		}
		return
	case streamcommons.BinanceStreamRESTDepth:
		symbolCap := strings.ToUpper(symbol)
		depth := new(jsonstructs.BinanceDepthREST)
		serr := json.Unmarshal(line, depth)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: depthrest BinanceDepthREST: %v", serr)
			return
		}
		formatted = make([]Result, len(depth.Asks)+len(depth.Bids))
		i := 0
		for _, order := range depth.Asks {
			fo := new(jsondef.BinanceRestDepth)
			fo.Symbol = symbolCap
			fo.Price, serr = strconv.ParseFloat(order[0], 64)
			if serr != nil {
				err = fmt.Errorf("FormatMessage: depthrest ask price: %v", serr)
				return
			}
			size, serr := strconv.ParseFloat(order[1], 64)
			if serr != nil {
				err = fmt.Errorf("FormatMessage: depthrest ask size: %v", serr)
				return
			}
			// Negative size for sell order
			fo.Size = -size
			mfo, serr := json.Marshal(fo)
			if serr != nil {
				err = fmt.Errorf("FormatMessage: depthrest ask BinanceDepth: %v", serr)
				return
			}
			formatted[i] = Result{
				Channel: channel,
				Message: mfo,
			}
			i++
		}
		for _, order := range depth.Bids {
			fo := new(jsondef.BinanceDepth)
			fo.Symbol = symbolCap
			fo.Price, serr = strconv.ParseFloat(order[0], 64)
			if serr != nil {
				err = fmt.Errorf("FormatMessage: depthrest bid price: %v", serr)
				return
			}
			fo.Size, serr = strconv.ParseFloat(order[1], 64)
			if serr != nil {
				err = fmt.Errorf("FormatMessage: depthrest bid size: %v", serr)
				return
			}
			mfo, serr := json.Marshal(fo)
			if serr != nil {
				err = fmt.Errorf("FormatMessage: depthrest bid BinanceDepth: %v", serr)
				return
			}
			formatted[i] = Result{
				Channel: channel,
				Message: mfo,
			}
			i++
		}
		return
	case streamcommons.BinanceStreamTrade:
		root := new(jsonstructs.BinanceReponseRoot)
		serr := json.Unmarshal(line, root)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: trade BinanceReponseRoot: %v", serr)
			return
		}
		trade := new(jsonstructs.BinanceTrade)
		serr = json.Unmarshal(root.Data, trade)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: trade BinanceTrade: %v", serr)
			return
		}
		ft := new(jsondef.BinanceTrade)
		ft.EventTime = strconv.FormatInt(trade.EventTime*int64(time.Millisecond), 10)
		ft.Timestamp = strconv.FormatInt(trade.TradeTime*int64(time.Millisecond), 10)
		ft.Symbol = trade.Symbol
		ft.Price, serr = strconv.ParseFloat(trade.Price, 64)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: trade price: %v", serr)
			return
		}
		ft.Size, serr = strconv.ParseFloat(trade.Quantity, 64)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: trade quantity: %v", serr)
			return
		}
		if trade.IsBuyerMarketMaker {
			// Buyer is the maker = seller is the taker
			// Negative size for sell order
			ft.Size = -ft.Size
		}
		ft.SellterOrderID = trade.SellerOrderID
		ft.BuyerOrderID = trade.BuyerOrderID
		ft.TradeID = trade.TradeID
		formatted = make([]Result, 1)
		mft, serr := json.Marshal(ft)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: trade BinanceTrade: %v", serr)
			return
		}
		formatted[0] = Result{
			Channel: channel,
			Message: mft,
		}
		return
	case streamcommons.BinanceStreamTicker:
		root := new(jsonstructs.BinanceReponseRoot)
		serr := json.Unmarshal(line, root)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: ticker BinanceReponseRoot: %v", serr)
			return
		}
		ticker := new(jsonstructs.BinanceTickerStream)
		serr = json.Unmarshal(root.Data, ticker)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: ticker BinanceTickerStream: %v", serr)
			return
		}
		ft := new(jsondef.BinanceTicker)
		ft.EventTime = strconv.FormatInt(ticker.EventTime*int64(time.Millisecond), 10)
		ft.Symbol = ticker.Symbol
		priceChange, serr := strconv.ParseFloat(ticker.PriceChange, 64)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: ticker priceChange: %v", serr)
			return
		}
		ft.PriceChange = priceChange
		priceChangePercent, serr := strconv.ParseFloat(ticker.PriceChanePercent, 64)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: ticker priceChangePercent: %v", serr)
			return
		}
		ft.PriceChanePercent = priceChangePercent
		weightedAveragePrice, serr := strconv.ParseFloat(ticker.WeightedAveragePrice, 64)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: ticker weightedAveragePrice: %v", serr)
			return
		}
		ft.WeightedAveragePrice = weightedAveragePrice
		firstTradePrice, serr := strconv.ParseFloat(ticker.FirstTradePrice, 64)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: ticker firstTradePrice: %v", serr)
			return
		}
		ft.FirstTradePrice = firstTradePrice
		lastPrice, serr := strconv.ParseFloat(ticker.LastPrice, 64)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: ticker lastPrice: %v", serr)
			return
		}
		ft.LastPrice = lastPrice
		lastQuantity, serr := strconv.ParseFloat(ticker.LastQuantity, 64)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: ticker lastQuantity: %v", serr)
			return
		}
		ft.LastQuantity = lastQuantity
		bestBidPrice, serr := strconv.ParseFloat(ticker.BestBidPrice, 64)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: ticker bestBidPrice: %v", serr)
			return
		}
		ft.BestBidPrice = bestBidPrice
		bestBidQuantity, serr := strconv.ParseFloat(ticker.BestBidQuantity, 64)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: ticker bestBidQuantity: %v", serr)
			return
		}
		ft.BestBidQuantity = bestBidQuantity
		bestAskPrice, serr := strconv.ParseFloat(ticker.BestAskPrice, 64)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: ticker bestAskPrice: %v", serr)
			return
		}
		ft.BestAskPrice = bestAskPrice
		bestAskQuantity, serr := strconv.ParseFloat(ticker.BestAskQuantity, 64)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: ticker bestAskQuantity: %v", serr)
			return
		}
		ft.BestAskQuantity = bestAskQuantity
		openPrice, serr := strconv.ParseFloat(ticker.OpenPrice, 64)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: ticker openPrice: %v", serr)
			return
		}
		ft.OpenPrice = openPrice
		highPrice, serr := strconv.ParseFloat(ticker.HighPrice, 64)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: ticker highPrice: %v", serr)
			return
		}
		ft.HighPrice = highPrice
		lowPrice, serr := strconv.ParseFloat(ticker.LowPrice, 64)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: ticker lowPrice: %v", serr)
			return
		}
		ft.LowPrice = lowPrice
		totalTradedBaseAssetVolume, serr := strconv.ParseFloat(ticker.TotalTradedBaseAssetVolume, 64)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: ticker totalTradedBaseAssetVolume: %v", serr)
			return
		}
		ft.TotalTradedBaseAssetVolume = totalTradedBaseAssetVolume
		totalTradedQuoteAssetVolume, serr := strconv.ParseFloat(ticker.TotalTradedQuoteAssetVolume, 64)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: ticker totalTradedQuoteAssetVolume: %v", serr)
			return
		}
		ft.TotalTradedQuoteAssetVolume = totalTradedQuoteAssetVolume
		ft.StatisticsOpenTime = strconv.FormatInt(ticker.StatisticsOpenTime*int64(time.Millisecond), 10)
		ft.StatisticsCloseTime = strconv.FormatInt(ticker.StatisticsCloseTime*int64(time.Millisecond), 10)
		ft.FirstTradeID = ticker.FirstTradeID
		ft.LastTradeID = ticker.LastTradeID
		ft.TotalNumberOfTrades = ticker.TotalNumberOfTrades
		formatted = make([]Result, 1)
		mft, serr := json.Marshal(ft)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: ticker BinanceTicker: %v", serr)
			return
		}
		formatted[0] = Result{
			Channel: channel,
			Message: mft,
		}
		return
	default:
		err = fmt.Errorf("FormatMessage: unsupported: %v", channel)
		return
	}
}

// IsSupported returns true if the given channel is supported by this formatter.
func (f *binanceFormatter) IsSupported(channel string) bool {
	_, stream, serr := streamcommons.BinanceDecomposeChannel(channel)
	if serr != nil {
		return false
	}
	return stream == streamcommons.BinanceStreamDepth ||
		stream == streamcommons.BinanceStreamTrade ||
		stream == streamcommons.BinanceStreamTicker ||
		stream == streamcommons.BinanceStreamRESTDepth
}
