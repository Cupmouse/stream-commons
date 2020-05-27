package csvlike

import (
	"encoding/json"
	"fmt"

	"github.com/exchangedataset/streamcommons/jsonstructs"
)

// BitmexFormatter formats message from bitmex
type BitmexFormatter struct {
}

// Format formats incoming message given and returns formatted strings
func (f *BitmexFormatter) Format(channel string, line []byte) (ret []string, err error) {
	subscribed := jsonstructs.BitmexSubscribe{}
	cerr := json.Unmarshal(line, &subscribed)
	if cerr != nil {
		err = fmt.Errorf("unmarshal failed: %s", cerr.Error())
		return
	}
	if subscribed.Success {
		// this is a response to subscription
		// return header row
		if channel == "orderBookL2" {
			ret = []string{HeaderBitmexOrderBookL2}
		} else if channel == "trade" {
			ret = []string{HeaderBitmexTrade}
		} else {
			err = fmt.Errorf("unsupported channel '%s' for csvlike formatting", channel)
		}
		return
	}

	if channel == "orderBookL2" {
		unmarshaled := jsonstructs.BitmexOrderBookL2{}
		cerr := json.Unmarshal(line, &unmarshaled)
		if cerr != nil {
			err = fmt.Errorf("unmarshal failed: %s", cerr.Error())
			return
		}

		ret = make([]string, len(unmarshaled.Data))
		for i, order := range unmarshaled.Data {
			size := int64(order.Size)
			if order.Side == "Sell" {
				// if side is sell, negate size
				size = -size
			}
			// pair, price, size
			ret[i] = fmt.Sprintf("%s,%d,%.8f,%d", order.Symbol, order.ID, order.Price, size)
		}
		return
	} else if channel == "trade" {
		unmarshaled := jsonstructs.BitmexTrade{}
		cerr := json.Unmarshal(line, &unmarshaled)
		if cerr != nil {
			err = fmt.Errorf("unmarshal failed: %s", cerr.Error())
			return
		}
		// Pair,Price,Size,Timestamp,TrdMatchID,TickDirection,GrossValue,HomeNotional,ForeignNotional
		ret = make([]string, len(unmarshaled.Data))
		for i, elem := range unmarshaled.Data {
			size := int64(elem.Size)
			if elem.Side == "Sell" {
				size = -size
			}
			grossValue := "null"
			if elem.GrossValue != nil {
				grossValue = fmt.Sprintf("%d", *elem.GrossValue)
			}
			homeNotional := "null"
			if elem.HomeNotional != nil {
				homeNotional = fmt.Sprintf("%.8f", *elem.HomeNotional)
			}
			foreignNotional := "null"
			if elem.ForeignNotional != nil {
				foreignNotional = fmt.Sprintf("%.8f", *elem.ForeignNotional)
			}
			ret[i] = fmt.Sprintf("%s,%.8f,%d,%s,%s,%s,%s,%s,%s", elem.Symbol, elem.Price, size, elem.Timestamp, elem.TradeMatchID, elem.TickDirection, grossValue, homeNotional, foreignNotional)
		}
		return
	}

	err = fmt.Errorf("channel '%s' if not supported for csvlike formatting", channel)
	return
}

// IsSupported returns true if given channel is supported to be formatted using this formatter
func (f *BitmexFormatter) IsSupported(channel string) bool {
	return channel == "orderBookL2" || channel == "trade"
}
