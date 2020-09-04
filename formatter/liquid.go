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

// liquidFormatter formats messages from Liquid to json format.
type liquidFormatter struct {
}

// FormatStart returns empty slice.
func (f *liquidFormatter) FormatStart(urlStr string) (formatted []StartReturn, err error) {
	return nil, nil
}

// FormatMessage formats raw messages from Liquid server to json format.
func (f *liquidFormatter) FormatMessage(channel string, line []byte) (formatted [][]byte, err error) {
	r := new(jsonstructs.LiquidMessageRoot)
	serr := json.Unmarshal(line, r)
	if serr != nil {
		err = fmt.Errorf("FormatMessage: root: %v", serr)
		return
	}
	if r.Event == jsonstructs.LiquidEventSubscriptionSucceeded {
		if strings.HasPrefix(channel, streamcommons.LiquidChannelPrefixLaddersCash) {
			return [][]byte{jsondef.TypeDefLiquidPriceLaddersCash}, nil
		} else if strings.HasPrefix(channel, streamcommons.LiquidChannelPrefixExecutionsCash) {
			return [][]byte{jsondef.TypeDefLiquidExecutionsCash}, nil
		} else {
			return nil, fmt.Errorf("FormatMessage: channel not supported: %v", channel)
		}
	}
	if strings.HasPrefix(channel, streamcommons.LiquidChannelPrefixLaddersCash) {
		// `true` is ask
		side := strings.HasSuffix(channel, "sell")
		channelPrefixStart := strings.LastIndexByte(channel, '_')
		if channelPrefixStart == -1 {
			return nil, fmt.Errorf("FormatMessage: price ladder no underscore in channel")
		}
		symbol := channel[len(streamcommons.LiquidChannelPrefixLaddersCash):channelPrefixStart]
		// Data for book is encoded in string
		dataStr := new(string)
		serr = json.Unmarshal(r.Data, dataStr)
		if serr != nil {
			return nil, fmt.Errorf("FormatMessage: price ladder dataStr: %v", serr)
		}
		orderbook := make([][]string, 0, 100)
		serr = json.Unmarshal([]byte(*dataStr), &orderbook)
		if serr != nil {
			return nil, fmt.Errorf("FormatMessage: price ladder orderbook: %v", serr)
		}
		formatted = make([][]byte, len(orderbook))
		for i, memOrder := range orderbook {
			price, serr := strconv.ParseFloat(memOrder[0], 64)
			if serr != nil {
				return nil, fmt.Errorf("FormatMessage: price: %v", serr)
			}
			quantity, serr := strconv.ParseFloat(memOrder[1], 64)
			if serr != nil {
				return nil, fmt.Errorf("FormatMessage: quantity: %v", serr)
			}
			order := new(jsondef.LiquidPriceLaddersCash)
			order.Symbol = symbol
			order.Price = price
			if side {
				order.Size = -quantity
			} else {
				order.Size = quantity
			}
			om, serr := json.Marshal(order)
			if serr != nil {
				return nil, fmt.Errorf("FormatMessage: order marshal: %v", serr)
			}
			formatted[i] = om
		}
		return formatted, nil
	} else if strings.HasPrefix(channel, streamcommons.LiquidChannelPrefixExecutionsCash) {
		execution := new(jsonstructs.LiquidExecution)
		dataStr := new(string)
		serr = json.Unmarshal(r.Data, dataStr)
		if serr != nil {
			return nil, fmt.Errorf("FormatMessage: execution dataStr: %v", serr)
		}
		serr = json.Unmarshal([]byte(*dataStr), execution)
		if serr != nil {
			return nil, fmt.Errorf("FormatMessage: execution: %v", serr)
		}
		lec := new(jsondef.LiquidExecutionsCash)
		createdAt := time.Unix(int64(execution.CreatedAt), 0)
		lec.CreatedAt = strconv.FormatInt(createdAt.UnixNano(), 10)
		lec.ID = execution.ID
		lec.Symbol = channel[len(streamcommons.LiquidChannelPrefixExecutionsCash):]
		lec.Price = execution.Price
		if execution.TakerSide == "sell" {
			lec.Size = -execution.Quantity
		} else {
			lec.Size = execution.Quantity
		}
		fm, serr := json.Marshal(lec)
		if serr != nil {
			return nil, fmt.Errorf("FormatMessage: formatted: %v", serr)
		}
		formatted = [][]byte{fm}
		return formatted, nil
	}
	return nil, fmt.Errorf("FormatMessage: line not supported")
}

func (f *liquidFormatter) IsSupported(channel string) bool {
	return (strings.HasPrefix(channel, streamcommons.LiquidChannelPrefixLaddersCash) && (strings.HasSuffix(channel, "buy")) || strings.HasSuffix(channel, "sell")) ||
		strings.HasPrefix(channel, streamcommons.LiquidChannelPrefixExecutionsCash)
}
