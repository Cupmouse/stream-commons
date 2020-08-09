package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/exchangedataset/streamcommons"
	"github.com/exchangedataset/streamcommons/formatter/json/jsondef"
	"github.com/exchangedataset/streamcommons/jsonstructs"
)

type binanceFormatter struct{}

func (f *binanceFormatter) Format(channel string, line []byte) (formatted [][]byte, err error) {
	subscribe := new(jsonstructs.BinanceSubscribe)
	serr := json.Unmarshal(line, subscribe)
	if serr != nil {
		err = fmt.Errorf("subscribe unmarshal: %v", serr)
		return
	}
	symbol, stream, serr := streamcommons.BinanceDecomposeChannel(channel)
	if serr != nil {
		err = serr
		return
	}
	if subscribe.Method != "SUBSCRIBE" {
		// Subscribe message
		formatted = make([][]byte, 1)
		if len(subscribe.Params) != 1 {
			err = errors.New("multiple channels in subscribe")
			return
		}
		if subscribe.Params[0] != channel {
			err = fmt.Errorf("channel differs: %v, expected: %v", subscribe.Params[0], channel)
			return
		}
		switch stream {
		case "depth":
			formatted[0] = jsondef.TypeDefBinanceDepth
		case "trade":
			formatted[0] = jsondef.TypeDefBinanceTrade
		case streamcommons.BinanceStreamRESTDepth:
			formatted[0] = jsondef.TypeDefBinanceRestDepth
		}
		return
	}
	switch stream {
	case "depth":
		root := new(jsonstructs.BinanceReponseRoot)
		serr := json.Unmarshal(line, root)
		if serr != nil {
			err = fmt.Errorf("root unmarshal: %v", serr)
			return
		}
		depth := new(jsonstructs.BinanceDepthStream)
		serr = json.Unmarshal(line, depth)
		if serr != nil {
			err = fmt.Errorf("depth unmarshal: %v", serr)
			return
		}
		formatted = make([][]byte, len(depth.Asks)+len(depth.Bids))
		i := 0
		for _, order := range depth.Asks {
			fo := new(jsondef.BinanceDepth)
			fo.Pair = depth.Symbol
			fo.Price, serr = strconv.ParseFloat(order[0], 64)
			if serr != nil {
				err = fmt.Errorf("Price ParseFloat: %v", serr)
				return
			}
			size, serr := strconv.ParseFloat(order[1], 64)
			if serr != nil {
				err = fmt.Errorf("Size ParseFloat: %v", serr)
				return
			}
			// Negative size for sell order
			fo.Size = -size
			mfo, serr := json.Marshal(fo)
			if serr != nil {
				err = fmt.Errorf("order marshal: %v", serr)
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
				err = fmt.Errorf("Price ParseFloat: %v", serr)
				return
			}
			fo.Size, serr = strconv.ParseFloat(order[1], 64)
			if serr != nil {
				err = fmt.Errorf("Size ParseFloat: %v", serr)
				return
			}
			mfo, serr := json.Marshal(fo)
			if serr != nil {
				err = fmt.Errorf("order marshal: %v", serr)
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
			err = fmt.Errorf("depth unmarshal: %v", serr)
			return
		}
		formatted = make([][]byte, len(depth.Asks)+len(depth.Bids))
		i := 0
		for _, order := range depth.Asks {
			fo := new(jsondef.BinanceDepth)
			fo.Pair = symbolCap
			fo.Price, serr = strconv.ParseFloat(order[0], 64)
			if serr != nil {
				err = fmt.Errorf("Price ParseFloat: %v", serr)
				return
			}
			size, serr := strconv.ParseFloat(order[1], 64)
			if serr != nil {
				err = fmt.Errorf("Size ParseFloat: %v", serr)
				return
			}
			// Negative size for sell order
			fo.Size = -size
			mfo, serr := json.Marshal(fo)
			if serr != nil {
				err = fmt.Errorf("order marshal: %v", serr)
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
				err = fmt.Errorf("Price ParseFloat: %v", serr)
				return
			}
			fo.Size, serr = strconv.ParseFloat(order[1], 64)
			if serr != nil {
				err = fmt.Errorf("Size ParseFloat: %v", serr)
				return
			}
			mfo, serr := json.Marshal(fo)
			if serr != nil {
				err = fmt.Errorf("order marshal: %v", serr)
				return
			}
			formatted[i] = mfo
			i++
		}
		return
	case "trade":
		root := new(jsonstructs.BinanceReponseRoot)
		serr := json.Unmarshal(line, root)
		if serr != nil {
			err = fmt.Errorf("root unmarshal: %v", serr)
			return
		}
		trade := new(jsonstructs.BinanceTrade)
		serr = json.Unmarshal(line, trade)
		if serr != nil {
			err = fmt.Errorf("trade unmarshal: %v", serr)
			return
		}
		ft := new(jsondef.BinanceTrade)
		ft.Eventtime = strconv.FormatInt(trade.EventTime*int64(time.Millisecond), 10)
		ft.Timestamp = strconv.FormatInt(trade.TradeTime*int64(time.Millisecond), 10)
		ft.Pair = trade.Symbol
		ft.Price, serr = strconv.ParseFloat(trade.Price, 64)
		if serr != nil {
			err = fmt.Errorf("Price ParseFloat: %v", serr)
			return
		}
		ft.Size, serr = strconv.ParseFloat(trade.Quantity, 64)
		if serr != nil {
			err = fmt.Errorf("Quantity ParseFloat: %v", serr)
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
			err = fmt.Errorf("trade marshal: %v", serr)
			return
		}
		formatted[0] = mft
		return
	default:
		err = fmt.Errorf("unsupported channel: %v", channel)
		return
	}
}

func (f *binanceFormatter) IsSupported(channel string) bool {
	_, stream, serr := streamcommons.BinanceDecomposeChannel(channel)
	if serr != nil {
		return false
	}
	return stream == "depth" || stream == "trade" ||
		stream == streamcommons.BinanceStreamRESTDepth
}
