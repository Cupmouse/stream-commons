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

func bitflyerParseTimestamp(timestamp string) (string, error) {
	timestampTime, serr := time.Parse(time.RFC3339Nano, timestamp)
	if serr != nil {
		return "", serr
	}
	return strconv.FormatInt(timestampTime.UnixNano(), 10), nil
}

// bitflyerFormatter formats raw input from bitflyer api into csvlike format.
type bitflyerFormatter struct {
}

// FormatStart returns empty slice.
func (f *bitflyerFormatter) FormatStart(urlStr string) ([]Result, error) {
	return make([]Result, 0), nil
}

func (f *bitflyerFormatter) formatBoard(channel string, messageRaw json.RawMessage) ([]Result, error) {
	var pair string
	if strings.HasPrefix(channel, "lightning_board_snapshot_") {
		pair = channel[len("lightning_board_snapshot_"):]
	} else {
		pair = channel[len("lightning_board_"):]
	}

	message := new(jsonstructs.BitflyerBoardParamsMessage)
	err := json.Unmarshal(messageRaw, message)
	if err != nil {
		return nil, fmt.Errorf("formatBoard: messageRaw: %v", err)
	}
	ret := make([]Result, len(message.Bids)+len(message.Asks))
	i := 0
	for _, ask := range message.Asks {
		// size == 0 if due to be removed
		marshaled, serr := json.Marshal(jsondef.BitflyerBoard{
			Symbol: pair,
			Price:  ask.Price,
			Side:   streamcommons.CommonFormatSell,
			Size:   ask.Size,
		})
		if serr != nil {
			return nil, fmt.Errorf("formatBoard: ask BitflyerBoard: %v", serr)
		}
		ret[i] = Result{
			Channel: channel,
			Message: marshaled,
		}
		i++
	}
	for _, bid := range message.Bids {
		marshaled, serr := json.Marshal(jsondef.BitflyerBoard{
			Symbol: pair,
			Price:  bid.Price,
			Side:   streamcommons.CommonFormatBuy,
			Size:   bid.Size,
		})
		if serr != nil {
			return nil, fmt.Errorf("formatBoard: bid BitflyerBoard: %v", serr)
		}
		ret[i] = Result{
			Channel: channel,
			Message: marshaled,
		}
		i++
	}
	return ret, nil
}

func (f *bitflyerFormatter) formatExecutions(channel string, messageRaw json.RawMessage) ([]Result, error) {
	// pair, price, size
	pair := channel[len("lightning_executions_"):]

	orders := make([]jsonstructs.BitflyerExecutionsParamMessageElement, 0, 10)
	err := json.Unmarshal(messageRaw, &orders)
	if err != nil {
		return nil, fmt.Errorf("formatExecutions: messageRaw: %v", err)
	}
	ret := make([]Result, len(orders))
	for i, element := range orders {
		var side string
		if element.Side == "ask" || element.Side == "SELL" {
			side = streamcommons.CommonFormatSell
		} else if element.Side == "bid" || element.Side == "BUY" {
			side = streamcommons.CommonFormatBuy
		} else {
			// for some reason, side could be empty, probably a bug of bitflyer api
			side = streamcommons.CommonFormatUnknown
		}
		marshaled, serr := json.Marshal(jsondef.BitflyerExecutions{
			Symbol: pair,
			Price:  element.Price,
			Side:   side,
			Size:   element.Size,
		})
		if serr != nil {
			return nil, fmt.Errorf("formatExecutions BitflyerExecutions: %v", serr)
		}
		ret[i] = Result{
			Channel: channel,
			Message: marshaled,
		}
	}
	return ret, nil
}

func (f *bitflyerFormatter) formatTicker(channel string, messageRaw json.RawMessage) ([]Result, error) {
	ticker := new(jsonstructs.BitflyerTickerParamsMessage)
	err := json.Unmarshal(messageRaw, ticker)
	if err != nil {
		return nil, fmt.Errorf("formatTicker: messageRaw: %v", err)
	}
	timestamp, serr := bitflyerParseTimestamp(ticker.Timestamp)
	if serr != nil {
		return nil, fmt.Errorf("formatTicker: timestamp: %v", serr)
	}

	marshaled, serr := json.Marshal(jsondef.BitflyerTicker{
		ProductCode:     ticker.ProductCode,
		Timestamp:       timestamp,
		TickID:          ticker.TickID,
		BestBid:         ticker.BestBid,
		BestAsk:         ticker.BestAsk,
		BestBidSize:     ticker.BestBidSize,
		BestAskSize:     ticker.BestAskSize,
		TotalBidDepth:   ticker.TotalBidDepth,
		TotalAskDepth:   ticker.TotalAskDepth,
		Ltp:             ticker.Ltp,
		Volume:          ticker.Volume,
		VolumeByProduct: ticker.VolumeByProduct,
	})
	if serr != nil {
		return nil, fmt.Errorf("formatTicker: BitflyerTicker: %v", serr)
	}
	return []Result{
		Result{
			Channel: channel,
			Message: marshaled,
		},
	}, nil
}

// FormatMessage formats message from bitflyer channel both given and returns formatted message
// keep in mind that multiple string will be returned
// error will be returned if channel is not supported to be formatted or
// message given is in invalid format
func (f *bitflyerFormatter) FormatMessage(channel string, line []byte) ([]Result, error) {
	// check if this message is a response to subscribe
	subscribe := jsonstructs.BitflyerSubscribed{}
	err := json.Unmarshal(line, &subscribe)
	if err != nil {
		return nil, fmt.Errorf("FormatMessage: line: %v", err)
	}
	if subscribe.Result {
		// an response for subscribe request
		if strings.HasPrefix(channel, "lightning_board_") {
			// lightning_board_snapshot will also return the same header
			return []Result{
				Result{
					Channel: channel,
					Message: jsondef.TypeDefBitflyerBoard,
				},
			}, nil
		} else if strings.HasPrefix(channel, "lightning_executions_") {
			return []Result{
				Result{
					Channel: channel,
					Message: jsondef.TypeDefBitflyerExecutions,
				},
			}, nil
		} else if strings.HasPrefix(channel, "lightning_ticker_") {
			return []Result{
				Result{
					Channel: channel,
					Message: jsondef.TypeDefBitflyerTicker,
				},
			}, nil
		} else {
			return nil, fmt.Errorf("csvlike unsupported: %s", channel)
		}
	} else {
		root := new(jsonstructs.BitflyerRoot)
		serr := json.Unmarshal(line, &root)
		if serr != nil {
			return nil, fmt.Errorf("FormatMessage: line: %v", err)
		}
		if strings.HasPrefix(channel, "lightning_board_") {
			return f.formatBoard(channel, root.Params.Message)
		} else if strings.HasPrefix(channel, "lightning_executions_") {
			return f.formatExecutions(channel, root.Params.Message)
		} else if strings.HasPrefix(channel, "lightning_ticker_") {
			return f.formatTicker(channel, root.Params.Message)
		} else {
			return nil, fmt.Errorf("csvlike unsupported: %s", channel)
		}
	}
}

// IsSupported returns true if message from given channel is supported to be formatted by this formatted
func (f *bitflyerFormatter) IsSupported(channel string) bool {
	return strings.HasPrefix(channel, "lightning_board_snapshot_") ||
		strings.HasPrefix(channel, "lightning_board_") ||
		strings.HasPrefix(channel, "lightning_executions_") ||
		strings.HasPrefix(channel, "lightning_ticker_")
}
