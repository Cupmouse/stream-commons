package jsonf

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/exchangedataset/streamcommons"
	"github.com/exchangedataset/streamcommons/formatter/jsonf/jsondef"
	"github.com/exchangedataset/streamcommons/jsonstructs"
)

// LiquidFormatter formats messages from Liquid to json format.
type LiquidFormatter struct {
}

// FormatStart returns empty slice.
func (f *LiquidFormatter) FormatStart(urlStr []byte) (formatted [][]byte, err error) {
	return make([][]byte, 0), nil
}

// FormatMessage formats raw messages from Liquid server to json format.
func (f *LiquidFormatter) FormatMessage(channel string, line []byte) (ret [][]byte, err error) {
	r := new(jsonstructs.LiquidMessageRoot)
	serr := json.Unmarshal(line, r)
	if serr != nil {
		err = fmt.Errorf("Format: root: %v", serr)
		return
	}
	if r.Event == jsonstructs.LiquidEventSubscriptionSucceeded {
		if strings.HasPrefix(channel, streamcommons.LiquidChannelPrefixLaddersCash) {
			return [][]byte{jsondef.TypeDefLiquidPriceLaddersCash}, nil
		} else if strings.HasPrefix(channel, streamcommons.LiquidChannelPrefixExecutionsCash) {
			return [][]byte{jsondef.TypeDefLiquidExecutionsCash}, nil
		} else {
			return nil, fmt.Errorf("Format: channel not supported: %v", channel)
		}
	}
	if strings.HasPrefix(channel, streamcommons.LiquidChannelPrefixLaddersCash) {
		// `true` is ask
		side := strings.HasSuffix(channel, "sell")
		orderbook := make([][]string, 0, 100)
		serr = json.Unmarshal(r.Data, &orderbook)
		if serr != nil {
			return nil, fmt.Errorf("Format: price ladder: %v", orderbook)
		}
		ret = make([][]byte, len(orderbook))
		for i, memOrder := range orderbook {
			price, serr := strconv.ParseFloat(memOrder[0], 64)
			if serr != nil {
				return nil, fmt.Errorf("Format: price: %v", serr)
			}
			quantity, serr := strconv.ParseFloat(memOrder[1], 64)
			if serr != nil {
				return nil, fmt.Errorf("Format: quantity: %v", serr)
			}
			order := new(jsondef.LiquidPriceLaddersCash)
			order.Price = price
			if side {
				order.Size = -quantity
			} else {
				order.Size = quantity
			}
			om, serr := json.Marshal(order)
			if serr != nil {
				return nil, fmt.Errorf("Format: order marshal: %v", serr)
			}
			ret[i] = om
		}
		return ret, nil
	} else if strings.HasPrefix(channel, streamcommons.LiquidChannelPrefixExecutionsCash) {
		execution := new(jsonstructs.LiquidExecution)
		serr = json.Unmarshal(r.Data, execution)
		if serr != nil {
			return nil, fmt.Errorf("Format: execution: %v", serr)
		}
		formatted := new(jsondef.LiquidExecutionsCash)
		createdAt := time.Unix(int64(execution.CreatedAt), 0)
		formatted.CreatedAt = createdAt.UnixNano()
		formatted.ID = execution.ID
		formatted.Pair = channel[len(streamcommons.LiquidChannelPrefixExecutionsCash):]
		formatted.Price = execution.Price
		if execution.TakerSide == "sell" {
			formatted.Size = -execution.Quantity
		} else {
			formatted.Size = execution.Quantity
		}
		fm, serr := json.Marshal(formatted)
		if serr != nil {
			return nil, fmt.Errorf("Format: formatted: %v", serr)
		}
		ret = [][]byte{fm}
	}
	return nil, fmt.Errorf("line not supported")
}
