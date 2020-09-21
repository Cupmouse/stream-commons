package streamcommons

import (
	"fmt"
	"strings"
)

// ChannelGroup is the group of channel used for categorizing the channel and its content to calculate the transfer cost.
type ChannelGroup int

const (
	// ChannelGroupOrderbook is the orderbook channel group.
	ChannelGroupOrderbook = ChannelGroup(iota)
	// ChannelGroupTrade is the orderbook channel group.
	ChannelGroupTrade
	// ChannelGroupOthers is the channel group for channels which is not included in any other channel groups.
	ChannelGroupOthers
)

// GetChannelGroup returns the channel group of the channel of the given exchange.
func GetChannelGroup(exchange string, channel string) (cg ChannelGroup, err error) {
	cg = ChannelGroupOthers
	switch exchange {
	case "bitmex":
		if strings.HasPrefix(channel, BitmexChannelOrderBookL2) {
			cg = ChannelGroupOrderbook
		} else if strings.HasPrefix(channel, BitmexChannelTrade) {
			cg = ChannelGroupTrade
		}
	case "bitfinex":
		if strings.HasPrefix(channel, BitfinexChannelPrefixBook) {
			cg = ChannelGroupOrderbook
		} else if strings.HasPrefix(channel, BitfinexChannelPrefixTrades) {
			cg = ChannelGroupTrade
		}
	case "bitflyer":
		if strings.HasPrefix(channel, BitflyerChannelPrefixLightningBoard) {
			// This includes board_snapshot channels
			cg = ChannelGroupOrderbook
		} else if strings.HasPrefix(channel, BitflyerChannelPrefixLightningExecutions) {
			cg = ChannelGroupTrade
		}
	case "binance":
		_, stream, serr := BinanceDecomposeChannel(channel)
		if serr != nil {
			err = fmt.Errorf("getChannelType: %v", serr)
			return
		}
		if stream == BinanceStreamTrade {
			cg = ChannelGroupTrade
		} else if stream == BinanceStreamDepth || stream == BinanceStreamRESTDepth {
			cg = ChannelGroupOrderbook
		}
	case "liquid":
		if strings.HasPrefix(channel, LiquidChannelPrefixLaddersCash) {
			cg = ChannelGroupOrderbook
		} else if strings.HasPrefix(channel, LiquidChannelPrefixExecutionsCash) {
			cg = ChannelGroupTrade
		}
	default:
		err = fmt.Errorf("getChannelType: exchange '%v' is not supported", exchange)
		return
	}
	return
}
