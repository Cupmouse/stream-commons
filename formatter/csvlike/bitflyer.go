package csvlike

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/exchangedataset/streamcommons/jsonstructs"
)

// BitflyerFormatter formats raw input from bitflyer api into csvlike format
type BitflyerFormatter struct {
}

func (f *BitflyerFormatter) formatBoard(channel string, line []byte) ([]string, error) {
	var pair string
	if strings.HasPrefix(channel, "lightning_board_snapshot_") {
		pair = channel[len("lightning_board_snapshot_"):]
	} else {
		pair = channel[len("lightning_board_"):]
	}

	unmarshaled := jsonstructs.BitflyerBoard{}
	err := json.Unmarshal(line, &unmarshaled)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed in formatBoard: %s", err.Error())
	}
	ret := make([]string, len(unmarshaled.Params.Message.Bids)+len(unmarshaled.Params.Message.Asks))
	i := 0
	for _, ask := range unmarshaled.Params.Message.Asks {
		ret[i] = fmt.Sprintf("%s,%.5f,%.8f", pair, ask.Price, -ask.Size)
		i++
	}
	for _, bid := range unmarshaled.Params.Message.Bids {
		// size == 0 if due to be removed
		ret[i] = fmt.Sprintf("%s,%.5f,%.8f", pair, bid.Price, bid.Size)
		i++
	}
	return ret, nil
}

func (f *BitflyerFormatter) formatExecutions(channel string, line []byte) ([]string, error) {
	// pair, price, size
	pair := channel[len("lightning_executions_"):]

	unmarshaled := jsonstructs.BitflyerExecutions{}
	err := json.Unmarshal(line, &unmarshaled)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed in formatExecutions: %s", err.Error())
	}
	ret := make([]string, len(unmarshaled.Params.Message))
	for i, element := range unmarshaled.Params.Message {
		if element.Side == "ask" {
			ret[i] = fmt.Sprintf("%s,%f.5,%f.8", pair, element.Price, -element.Size)
		} else {
			// for some reason, size could be empty, probably a bug of bitflyer api
			// this also includes buy side
			ret[i] = fmt.Sprintf("%s,%f.5,%f.8", pair, element.Price, element.Size)
		}
	}
	return ret, nil
}

// Format formats message from bitflyer channel both given and returns formatted message
// keep in mind that multiple string will be returned
// error will be returned if channel is not supported to be formatted or
// message given is in invalid format
func (f *BitflyerFormatter) Format(channel string, line []byte) ([]string, error) {
	// check if this message is a response to subscribe
	subscribe := jsonstructs.BitflyerSubscribed{}
	err := json.Unmarshal(line, &subscribe)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %s", err.Error())
	}
	if subscribe.Result {
		// an response for subscribe request
		if strings.HasPrefix(channel, "lightning_board_") {
			// lightning_board_snapshot will also return the same header
			return []string{HeaderOrderBook}, nil
		} else if strings.HasPrefix(channel, "lightning_executions_") {
			return []string{HeaderTrade}, nil
		} else {
			return nil, fmt.Errorf("unsupported channel for csvlike formatting: %s", channel)
		}
	} else {
		if strings.HasPrefix(channel, "lightning_board_") {
			return f.formatBoard(channel, line)
		} else if strings.HasPrefix(channel, "lightning_executions_") {
			return f.formatExecutions(channel, line)
		} else {
			return nil, fmt.Errorf("unsupported channel for csvlike formatting: %s", channel)
		}
	}
}

// IsSupported returns true if message from given channel is supported to be formatted by this formatted
func (f *BitflyerFormatter) IsSupported(channel string) bool {
	return strings.HasPrefix(channel, "lightning_board_snapshot_") ||
		strings.HasPrefix(channel, "lightning_board_") ||
		strings.HasPrefix(channel, "lightning_executions_")
}
