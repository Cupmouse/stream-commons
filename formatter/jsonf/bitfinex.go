package jsonf

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/exchangedataset/streamcommons/formatter/jsonf/jsondef"
	"github.com/exchangedataset/streamcommons/jsonstructs"
)

// BitfinexFormatter formats raw input from bitfinex api
type BitfinexFormatter struct{}

// FormatStart returns empty slice.
func (f *BitfinexFormatter) FormatStart(urlStr string) ([][]byte, error) {
	return make([][]byte, 0), nil
}

func (f *BitfinexFormatter) formatBook(channel string, line []byte) ([][]byte, error) {
	pair := channel[len("book_"):]

	unmarshaled := jsonstructs.BitfinexBook{}
	err := json.Unmarshal(line, &unmarshaled)
	if err != nil {
		return nil, fmt.Errorf("formatBook: line: %v", err)
	}
	var orders []jsonstructs.BitfinexBookOrder

	// if there is only one order in this message, it would get flattened
	// there is also heartbeat message where unmarshaled[1] is constant "hb"
	// bad api design by bitfinex
	switch unmarshaled[1].(type) {
	case string:
		// heartbeat
		return nil, nil
	case []interface{}:
		ordersInterf := unmarshaled[1].([]interface{})
		switch ordersInterf[0].(type) {
		case float64:
			// only one order
			orders = []jsonstructs.BitfinexBookOrder{ordersInterf}
			break
		case []interface{}:
			orders = make([]jsonstructs.BitfinexBookOrder, len(ordersInterf))
			for i, orderInterf := range ordersInterf {
				orders[i] = orderInterf.(jsonstructs.BitfinexBookOrder)
			}
		default:
			return nil, fmt.Errorf("formatBook: invalid order type: %s", reflect.TypeOf(ordersInterf[0]))
		}

		ret := make([][]byte, len(orders))
		// TODO adding funding pair support
		// normal trade pair
		for i, order := range orders {
			// order[0] = price, order[1] = count, order[2] = +-amount
			price := order[0].(float64)
			count := int64(order[1].(float64))
			size := order[2].(float64)
			if count == 0 {
				size = 0
			}

			marshaled, serr := json.Marshal(jsondef.BitfinexBook{
				Pair:  pair,
				Price: price,
				Count: count,
				Size:  size,
			})
			if serr != nil {
				return nil, fmt.Errorf("formatBook: BitfinexBook: %v", serr)
			}
			ret[i] = marshaled
		}
		return ret, nil
	default:
		return nil, fmt.Errorf("formatBook: invalid payload type: %s", reflect.TypeOf(unmarshaled[1]))
	}
}

func (f *BitfinexFormatter) formatTrades(channel string, line []byte) (formatted [][]byte, err error) {
	pair := channel[len("trades_"):]

	unmarshal := jsonstructs.BitfinexTrades{}
	err = json.Unmarshal(line, &unmarshal)
	if err != nil {
		return nil, fmt.Errorf("formatTrades: BitfinexTrades: %s", err)
	}

	var ordersInterf []interface{}
	switch unmarshal[1].(type) {
	case string:
		if unmarshal[1].(string) == "hb" {
			// heatbeat, ignore
			return nil, nil
		}
		teTu := unmarshal[1].(string)

		if teTu == "tu" {
			// bitfinex blog says tu includes tradeid, but it looks like te also has it
			return
		}
		ordersInterf = unmarshal[2].([]interface{})
		break
	case []interface{}:
		// first message does not have tu,te
		ordersInterf = unmarshal[1].([]interface{})
		break
	default:
		return nil, fmt.Errorf("formatTrades: first element invalid type: %s", reflect.TypeOf(unmarshal[1]))
	}

	var orders []jsonstructs.BitfinexTradesElement
	switch ordersInterf[0].(type) {
	case float64:
		// only one trade element
		orders = []jsonstructs.BitfinexTradesElement{ordersInterf}
		break
	case []interface{}:
		orders = make([]jsonstructs.BitfinexTradesElement, len(ordersInterf))
		for i, orderInterf := range ordersInterf {
			orders[i] = orderInterf.(jsonstructs.BitfinexTradesElement)
		}
		break
	default:
		err = fmt.Errorf("formatTrades: invalid order type: %s", reflect.TypeOf(ordersInterf[0]))
		return
	}

	formatted = make([][]byte, len(orders))
	for i, order := range orders {
		// orderID, millisectimestamp, amount, +-price
		orderID := int64(order[0].(float64))
		millisecTimestamp := uint64(order[1].(float64))
		size := order[2].(float64)
		price := order[3].(float64)

		// convert timestamp into nanosec
		timestamp := fmt.Sprintf("%d", millisecTimestamp*1000000)

		var marshaled []byte
		marshaled, err = json.Marshal(jsondef.BitfinexTrades{
			Pair:      pair,
			OrderID:   orderID,
			Timestamp: timestamp,
			Price:     price,
			Size:      size,
		})
		if err != nil {
			return
		}
		formatted[i] = marshaled
	}

	return
}

// FormatMessage formats line from channel given and returns an array of them
func (f *BitfinexFormatter) FormatMessage(channel string, line []byte) (formatted [][]byte, err error) {
	subscribe := jsonstructs.BitfinexSubscribed{}
	err = json.Unmarshal(line, &subscribe)
	if err == nil && subscribe.Event == "subscribed" {
		// an response for subscribe request
		if strings.HasPrefix(channel, "book_") {
			formatted = [][]byte{jsondef.TypeDefBitfinexBook}
			return
		} else if strings.HasPrefix(channel, "trades_") {
			formatted = [][]byte{jsondef.TypeDefBitfinexTrades}
			return
		} else {
			err = fmt.Errorf("FormatMessage: json unsupported channel: %s", channel)
			return
		}
	} else {
		if strings.HasPrefix(channel, "book_") {
			formatted, err = f.formatBook(channel, line)
			return
		} else if strings.HasPrefix(channel, "trades_") {
			formatted, err = f.formatTrades(channel, line)
			return
		} else {
			err = fmt.Errorf("FormatMessage: json unsupported channel: %s", channel)
			return
		}
	}
}

// IsSupported returns true if specified channel is supported to be formatted using this formatter
func (f *BitfinexFormatter) IsSupported(channel string) bool {
	return strings.HasPrefix(channel, "book_") ||
		strings.HasPrefix(channel, "trades_")
}
