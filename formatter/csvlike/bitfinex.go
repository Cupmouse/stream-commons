package csvlike

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/exchangedataset/streamcommons/jsonstructs"
)

// BitfinexFormatter formats raw input from bitfinex api
type BitfinexFormatter struct {
}

func (f *BitfinexFormatter) formatBook(channel string, line []byte) ([]string, error) {
	pair := channel[len("book_"):]

	unmarshaled := jsonstructs.BitfinexBook{}
	err := json.Unmarshal(line, &unmarshaled)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %s", err.Error())
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
			return nil, fmt.Errorf("invalid type for order: %s", reflect.TypeOf(ordersInterf[0]))
		}

		// order[0] = price, order[1] = count, order[2] = +-amount
		ret := make([]string, len(orders))
		// TODO adding funding pair support
		// normal trade pair
		for i, order := range orders {
			// pair, price, size
			ret[i] = fmt.Sprintf("%s,%.8f,%.8f", pair, order[0], order[2])
		}
		return ret, nil
	default:
		return nil, fmt.Errorf("invalid type for payload: %s", reflect.TypeOf(unmarshaled[1]))
	}
}

func (f *BitfinexFormatter) formatTrades(channel string, line []byte) ([]string, error) {
	pair := channel[len("trades_"):]

	unmarshal := jsonstructs.BitfinexTrades{}
	err := json.Unmarshal(line, &unmarshal)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %s", err.Error())
	}
	switch unmarshal[1].(type) {
	case string:
		// heartbeat
		return nil, nil
	case []interface{}:
		ordersInterf := unmarshal[1].([]interface{})

		var orders []jsonstructs.BitfinexTradesElement
		switch ordersInterf[0].(type) {
		case float64:
			// only one trade element
			orders = []jsonstructs.BitfinexTradesElement{ordersInterf}
			break
		case []interface{}:
			for i, orderInterf := range ordersInterf[0].([]interface{}) {
				orders[i] = orderInterf.(jsonstructs.BitfinexTradesElement)
			}
			break
		default:
			return nil, fmt.Errorf("invalid type for order: %s", reflect.TypeOf(ordersInterf[0]))
		}

		ret := make([]string, len(orders))
		for i, order := range orders {
			// pair, orderID, millisectimestamp, price, rate
			ret[i] = fmt.Sprintf("%s,%d,%d,%.8f,%.8f", pair, order[0], order[1], order[3], order[2])
		}

		return ret, nil
	default:
		return nil, fmt.Errorf("invalid type for payload: %s", reflect.TypeOf(unmarshal[1]))
	}
}

// Format formats line from channel given and returns an array of them
func (f *BitfinexFormatter) Format(channel string, line []byte) ([]string, error) {
	subscribe := jsonstructs.BitfinexSubscribed{}
	err := json.Unmarshal(line, &subscribe)
	if err == nil && subscribe.Event == "subscribed" {
		// an response for subscribe request
		if strings.HasPrefix(channel, "book_") {
			return []string{HeaderOrderBook}, nil
		} else if strings.HasPrefix(channel, "trades_") {
			return []string{HeaderTrade}, nil
		} else {
			return nil, fmt.Errorf("unsupported channel for csvlike formatting: %s", channel)
		}
	} else {
		if strings.HasPrefix(channel, "book_") {
			return f.formatBook(channel, line)
		} else if strings.HasPrefix(channel, "trades_") {
			return f.formatTrades(channel, line)
		} else {
			return nil, fmt.Errorf("unsupported channel for csvlike formatting: %s", channel)
		}
	}
}

// IsSupported returns true if specified channel is supported to be formatted using this formatter
func (f *BitfinexFormatter) IsSupported(channel string) bool {
	return strings.HasPrefix(channel, "book_") ||
		strings.HasPrefix(channel, "trades_")
}
