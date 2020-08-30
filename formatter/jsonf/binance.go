package jsonf

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/exchangedataset/streamcommons"
	"github.com/exchangedataset/streamcommons/formatter/jsonf/jsondef"
	"github.com/exchangedataset/streamcommons/jsonstructs"
)

// BinanceFormatter is json formatter for Binance.
type BinanceFormatter struct{}

// FormatStart formats start line (URL) and returns the array of known subscribed channel in case the server won't
// tell the client what channels are successfully subscribed.
func (f *BinanceFormatter) FormatStart(urlStr string) (formatted [][]byte, err error) {
	u, serr := url.Parse(string(urlStr))
	if serr != nil {
		return nil, fmt.Errorf("FormatStart: %v", serr)
	}
	q := u.Query()
	streams := q.Get("streams")
	channels := strings.Split("/", streams)
	formatted = make([][]byte, len(channels))
	for i, ch := range channels {
		_, stream, serr := streamcommons.BinanceDecomposeChannel(ch)
		if serr != nil {
			err = fmt.Errorf("FormatStart: %v", serr)
			return
		}
		switch stream {
		case streamcommons.BinanceStreamDepth:
			formatted[i] = jsondef.TypeDefBinanceDepth
		case streamcommons.BinanceStreamTrade:
			formatted[i] = jsondef.TypeDefBinanceTrade
		case streamcommons.BinanceStreamRESTDepth:
			formatted[i] = jsondef.TypeDefBinanceRestDepth
		default:
			err = fmt.Errorf("FormatStart: channel not supported: %s", ch)
			return
		}
	}
	return formatted, nil
}

// FormatMessage formats messages from server.
func (f *BinanceFormatter) FormatMessage(channel string, line []byte) (formatted [][]byte, err error) {
	symbol, stream, serr := streamcommons.BinanceDecomposeChannel(channel)
	if serr != nil {
		err = serr
		return
	}
	// subscribe := new(jsonstructs.BinanceSubscribe)
	// serr := json.Unmarshal(line, subscribe)
	// if serr != nil {
	// 	err = fmt.Errorf("FormatMessage: line: %v", serr)
	// 	return
	// }
	// if subscribe.Method != "SUBSCRIBE" {
	// 	// Subscribe message
	// 	formatted = make([][]byte, 1)
	// 	if len(subscribe.Params) != 1 {
	// 		err = errors.New("FormatMessage: multiple channels in subscribe")
	// 		return
	// 	}
	// 	if subscribe.Params[0] != channel {
	// 		err = fmt.Errorf("FormatMessage: channel differs: %v, expected: %v", subscribe.Params[0], channel)
	// 		return
	// 	}
	// 	switch stream {
	// 	case "depth":
	// 		formatted[0] = jsondef.TypeDefBinanceDepth
	// 	case "trade":
	// 		formatted[0] = jsondef.TypeDefBinanceTrade
	// 	case streamcommons.BinanceStreamRESTDepth:
	// 		formatted[0] = jsondef.TypeDefBinanceRestDepth
	// 	default:
	// 		err = fmt.Errorf("FormatMessage: channel not supported: %s", )
	// 	}
	// 	return
	// }
	switch stream {
	case streamcommons.BinanceStreamDepth:
		root := new(jsonstructs.BinanceReponseRoot)
		serr := json.Unmarshal(line, root)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: BinanceReponseRoot: %v", serr)
			return
		}
		depth := new(jsonstructs.BinanceDepthStream)
		serr = json.Unmarshal(line, depth)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: BinanceDepthStream: %v", serr)
			return
		}
		formatted = make([][]byte, len(depth.Asks)+len(depth.Bids))
		i := 0
		for _, order := range depth.Asks {
			fo := new(jsondef.BinanceDepth)
			fo.Pair = depth.Symbol
			fo.Price, serr = strconv.ParseFloat(order[0], 64)
			if serr != nil {
				err = fmt.Errorf("FormatMessage: price: %v", serr)
				return
			}
			size, serr := strconv.ParseFloat(order[1], 64)
			if serr != nil {
				err = fmt.Errorf("FormatMessage: size: %v", serr)
				return
			}
			// Negative size for sell order
			fo.Size = -size
			mfo, serr := json.Marshal(fo)
			if serr != nil {
				err = fmt.Errorf("FormatMessage: BinanceDepth: %v", serr)
				return
			}
			formatted[i] = mfo
			i++
		}
		for _, order := range depth.Bids {
			fo := new(jsondef.BinanceDepth)
			fo.Pair = depth.Symbol
			fo.Price, serr = strconv.ParseFloat(order[0], 64)
			if serr != nil {
				err = fmt.Errorf("FormatMessage: price: %v", serr)
				return
			}
			fo.Size, serr = strconv.ParseFloat(order[1], 64)
			if serr != nil {
				err = fmt.Errorf("FormatMessage: size: %v", serr)
				return
			}
			mfo, serr := json.Marshal(fo)
			if serr != nil {
				err = fmt.Errorf("FormatMessage: order marshal: %v", serr)
				return
			}
			formatted[i] = mfo
			i++
		}
		return
	case streamcommons.BinanceStreamRESTDepth:
		symbolCap := strings.ToUpper(symbol)
		depth := new(jsonstructs.BinanceDepthREST)
		serr := json.Unmarshal(line, depth)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: depth unmarshal: %v", serr)
			return
		}
		formatted = make([][]byte, len(depth.Asks)+len(depth.Bids))
		i := 0
		for _, order := range depth.Asks {
			fo := new(jsondef.BinanceDepth)
			fo.Pair = symbolCap
			fo.Price, serr = strconv.ParseFloat(order[0], 64)
			if serr != nil {
				err = fmt.Errorf("FormatMessage: price: %v", serr)
				return
			}
			size, serr := strconv.ParseFloat(order[1], 64)
			if serr != nil {
				err = fmt.Errorf("FormatMessage: size: %v", serr)
				return
			}
			// Negative size for sell order
			fo.Size = -size
			mfo, serr := json.Marshal(fo)
			if serr != nil {
				err = fmt.Errorf("FormatMessage: ask BinanceDepth: %v", serr)
				return
			}
			formatted[i] = mfo
			i++
		}
		for _, order := range depth.Bids {
			fo := new(jsondef.BinanceDepth)
			fo.Pair = symbolCap
			fo.Price, serr = strconv.ParseFloat(order[0], 64)
			if serr != nil {
				err = fmt.Errorf("FormatMessage: price: %v", serr)
				return
			}
			fo.Size, serr = strconv.ParseFloat(order[1], 64)
			if serr != nil {
				err = fmt.Errorf("FormatMessage: size: %v", serr)
				return
			}
			mfo, serr := json.Marshal(fo)
			if serr != nil {
				err = fmt.Errorf("FormatMessage: bid BinanceDepth: %v", serr)
				return
			}
			formatted[i] = mfo
			i++
		}
		return
	case streamcommons.BinanceStreamTrade:
		root := new(jsonstructs.BinanceReponseRoot)
		serr := json.Unmarshal(line, root)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: line: %v", serr)
			return
		}
		trade := new(jsonstructs.BinanceTrade)
		serr = json.Unmarshal(line, trade)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: BinanceTrade: %v", serr)
			return
		}
		ft := new(jsondef.BinanceTrade)
		ft.Eventtime = strconv.FormatInt(trade.EventTime*int64(time.Millisecond), 10)
		ft.Timestamp = strconv.FormatInt(trade.TradeTime*int64(time.Millisecond), 10)
		ft.Pair = trade.Symbol
		ft.Price, serr = strconv.ParseFloat(trade.Price, 64)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: price: %v", serr)
			return
		}
		ft.Size, serr = strconv.ParseFloat(trade.Quantity, 64)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: quantity: %v", serr)
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
		mft, serr := json.Marshal(ft)
		if serr != nil {
			err = fmt.Errorf("FormatMessage: BinanceTrade: %v", serr)
			return
		}
		formatted[0] = mft
		return
	default:
		err = fmt.Errorf("FormatMessage: unsupported: %v", channel)
		return
	}
}

// IsSupported returns true if the given channel is supported by this formatter.
func (f *BinanceFormatter) IsSupported(channel string) bool {
	_, stream, serr := streamcommons.BinanceDecomposeChannel(channel)
	if serr != nil {
		return false
	}
	return stream == streamcommons.BinanceStreamDepth ||
		stream == streamcommons.BinanceStreamTrade ||
		stream == streamcommons.BinanceStreamRESTDepth
}
