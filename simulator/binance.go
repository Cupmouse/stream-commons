package simulator

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/exchangedataset/streamcommons"
	"github.com/exchangedataset/streamcommons/jsonstructs"
)

type binanceOrderbook struct {
	asks map[float64]float64
	bids map[float64]float64
	// True if immediate last operation was to construct the initial state
	isLastSnapshot bool
	// Last received FinalUpdateID to check for missing messages
	lastFinalUpdateID int64
	// Orderbook differences waiting to be applied for the arrival of a REST message
	// nil if a REST message has already been received
	differences []*jsonstructs.BinanceDepthStream
}

type binanceSimulator struct {
	filterChannel map[string]bool
	// IDs and its channel a client sent a subscription to
	idCh map[int]string
	// Slice of channels a client subscribed to and a server agreed
	subscribed []string
	// map[symbol]orderbook
	// Note: symbol is lower-cased one
	orderBooks map[string]*binanceOrderbook
}

func (s *binanceSimulator) ProcessStart(line []byte) error {
	u, serr := url.Parse(string(line))
	if serr != nil {
		return serr
	}
	query := u.Query()
	channels := strings.Split(query.Get("streams"), "/")
	for _, ch := range channels {
		if s.filterChannel != nil {
			_, ok := s.filterChannel[ch]
			if !ok {
				continue
			}
		}
		// Add this channel to successfully subscribed channel list
		s.subscribed = append(s.subscribed, ch)
		symbol, stream, serr := streamcommons.BinanceDecomposeChannel(ch)
		if serr != nil {
			return serr
		}
		if stream == "depth@100ms" {
			// Create new orderbook in memory
			if _, ok := s.orderBooks[symbol]; ok {
				return errors.New("received subscribe confirmation twice")
			}
			orderbook := new(binanceOrderbook)
			orderbook.asks = make(map[float64]float64, 10000)
			orderbook.bids = make(map[float64]float64, 10000)
			// Create a slice to store difference messages before a REST message arrives
			orderbook.differences = make([]*jsonstructs.BinanceDepthStream, 0, 1000)
			s.orderBooks[symbol] = orderbook
		}
	}
	return nil
}

func (s *binanceSimulator) ProcessSend(line []byte) (channel string, err error) {
	channel = streamcommons.ChannelUnknown
	sub := new(jsonstructs.BinanceSubscribe)
	serr := json.Unmarshal(line, sub)
	if serr != nil {
		err = fmt.Errorf("subscribe unmarshal: %v", serr)
		return
	}
	if len(sub.Params) != 1 {
		// Subscription have to be done one at the time
		err = errors.New("len(subscribe.Params) != 1")
		return
	}
	if sub.ID == 0 {
		err = errors.New("use of 0 as subscription id")
	}
	channel = sub.Params[0]
	s.idCh[sub.ID] = channel
	return
}

func binanceProcessSide(asks [][]string, m map[float64]float64) (err error) {
	for _, order := range asks {
		price, serr := strconv.ParseFloat(order[0], 64)
		if serr != nil {
			err = fmt.Errorf("price ParseFloat: %v", serr)
			return
		}
		quantity, serr := strconv.ParseFloat(order[1], 64)
		if serr != nil {
			err = fmt.Errorf("quantity ParseFloat: %v", serr)
			return
		}
		if quantity == 0 {
			// Remove order from the book
			delete(m, price)
		} else {
			// Update order
			m[price] = quantity
		}
	}
	return nil
}

func (s *binanceSimulator) processMessageDepth(channel string, depth *jsonstructs.BinanceDepthStream) (err error) {
	symbol, _, serr := streamcommons.BinanceDecomposeChannel(channel)
	if serr != nil {
		return serr
	}
	orderbook := s.orderBooks[symbol]
	if orderbook.isLastSnapshot {
		// First event should have this traits
		if depth.FirstUpdateID > orderbook.lastFinalUpdateID+1 ||
			depth.FinalUpdateID < orderbook.lastFinalUpdateID+1 {
			err = fmt.Errorf("first difference's updateID out of range")
			return
		}
		orderbook.isLastSnapshot = false
	} else {
		if orderbook.lastFinalUpdateID+1 != depth.FirstUpdateID {
			// There are missing messages that haven't been received
			err = fmt.Errorf("missing messages detected")
			return
		}
	}
	err = binanceProcessSide(depth.Asks, orderbook.asks)
	if err != nil {
		return
	}
	err = binanceProcessSide(depth.Bids, orderbook.bids)
	if err != nil {
		return
	}
	orderbook.lastFinalUpdateID = depth.FinalUpdateID
	return nil
}

func (s *binanceSimulator) ProcessMessageWebSocket(line []byte) (channel string, err error) {
	channel = streamcommons.ChannelUnknown

	root := new(jsonstructs.BinanceReponseRoot)
	err = json.Unmarshal(line, root)
	if err != nil {
		return
	}
	channel = root.Stream
	symbol, stream, serr := streamcommons.BinanceDecomposeChannel(channel)
	if serr != nil {
		err = serr
		return
	}
	if stream == "depth@100ms" {
		depth := new(jsonstructs.BinanceDepthStream)
		serr := json.Unmarshal(root.Data, depth)
		if serr != nil {
			err = fmt.Errorf("depth message unmarshal: %v", serr)
			return
		}
		orderbook := s.orderBooks[symbol]
		if orderbook.lastFinalUpdateID != 0 {
			// Already received a REST message
			err = s.processMessageDepth(channel, depth)
			return
		}
		if len(orderbook.differences) > 100 {
			err = fmt.Errorf("too much stored difference: %v", symbol)
			return
		}
		// Store this message into an slice
		orderbook.differences = append(orderbook.differences, depth)
	}
	// Ignore other channels
	return
}

func (s *binanceSimulator) ProcessMessageChannelKnown(channel string, line []byte) (err error) {
	symbol, stream, serr := streamcommons.BinanceDecomposeChannel(channel)
	if serr != nil {
		return serr
	}
	if stream == streamcommons.BinanceStreamRESTDepth {
		// REST depth message
		depthRest := new(jsonstructs.BinanceDepthREST)
		serr := json.Unmarshal(line, depthRest)
		if serr != nil {
			return fmt.Errorf("depth unmarshal: %v", serr)
		}
		if depthRest.LastUpdateID == 0 {
			return errors.New("depth unmarshal: LastUpdateID == 0, probably not a depth message")
		}
		// Set the initial state
		orderbook := s.orderBooks[symbol]
		if orderbook.lastFinalUpdateID != 0 {
			err = errors.New("received REST twice")
			return
		}
		orderbook.isLastSnapshot = true
		orderbook.lastFinalUpdateID = depthRest.LastUpdateID
		err = binanceProcessSide(depthRest.Asks, orderbook.asks)
		if err != nil {
			return
		}
		err = binanceProcessSide(depthRest.Bids, orderbook.bids)
		if err != nil {
			return
		}
		// Apply orderbook differences previously received via WebSocket
		differences := orderbook.differences
		// Drop unneccesary stored messages
		i := 0
		for ; i < len(differences) && differences[i].FinalUpdateID <= depthRest.LastUpdateID; i++ {
		}
		if i == len(differences) {
			// No messages that should be applied immediately are stored
			return
		}
		// Apply all differences stored
		for ; i < len(differences); i++ {
			serr := s.processMessageDepth(symbol+"@depth@100ms", differences[i])
			if serr != nil {
				return fmt.Errorf("apply depth: %v", serr)
			}
		}
		// To free memory space
		orderbook.differences = nil
		return
	}
	wsChannel, serr := s.ProcessMessageWebSocket(line)
	if serr != nil {
		return serr
	}
	if wsChannel != channel {
		return fmt.Errorf("channel differs: %v expected: %v", wsChannel, channel)
	}
	return
}

func (s *binanceSimulator) ProcessState(channel string, line []byte) (err error) {
	if channel == streamcommons.StateChannelSubscribed {
		subscribed := make([]string, 0, 100)
		// Read subscribed state
		serr := json.Unmarshal(line, &subscribed)
		if serr != nil {
			err = fmt.Errorf("subscribed state unmarshal: %v", serr)
			return
		}
		if s.filterChannel == nil {
			// Filter is disabled, add all
			s.subscribed = subscribed
		} else {
			// Apply filter to it
			for _, stateChannel := range subscribed {
				_, ok := s.filterChannel[stateChannel]
				if ok {
					s.subscribed = append(s.subscribed, stateChannel)
				}
			}
		}
		return
	}
	symbol, stream, serr := streamcommons.BinanceDecomposeChannel(channel)
	if serr != nil {
		return
	}
	switch stream {
	case "depth@100ms":
		depthState := make(map[string]*binanceOrderbook)
		serr := json.Unmarshal(line, &depthState)
		if serr != nil {
			err = fmt.Errorf("depth state unmarshal: %v", serr)
			return
		}
		if s.filterChannel == nil {
			// Filter disabled
			s.orderBooks = depthState
		} else {
			// Apply filter
			for _, subChannel := range s.subscribed {
				_, ok := s.filterChannel[subChannel]
				if ok {
					s.orderBooks[symbol] = depthState[symbol]
				}
			}
		}
		return
	default:
		err = fmt.Errorf("unknown stream name: %v", stream)
		return
	}
}

func (s *binanceSimulator) TakeStateSnapshot() (snapshot []Snapshot, err error) {
	if s.filterChannel != nil {
		// If channel filtering is enabled, this should not be called
		err = errors.New("channel filter is enabled")
		return
	}
	snapshot = make([]Snapshot, 0, 100)
	// Take a snapshot of subscribed channel list
	subMarshaled, serr := json.Marshal(s.subscribed)
	if serr != nil {
		err = fmt.Errorf("subscribed marshal: %v", serr)
		return
	}
	snapshot = append(snapshot, Snapshot{
		Channel:  streamcommons.StateChannelSubscribed,
		Snapshot: subMarshaled,
	})
	// Take snapshots of orderbooks
	for symbol, orderbook := range s.orderBooks {
		depthRest := new(jsonstructs.BinanceDepthREST)
		depthRest.LastUpdateID = orderbook.lastFinalUpdateID
		depthRest.Asks = make([][]string, len(depthRest.Asks))

		obMarshaled, serr := json.Marshal(orderbook)
		if serr != nil {
			err = fmt.Errorf("orderbook marshal: %v", serr)
			return
		}
		snapshot = append(snapshot, Snapshot{
			Channel:  symbol + "@" + streamcommons.BinanceStreamRESTDepth,
			Snapshot: obMarshaled,
		})
	}
	return
}

func (s *binanceSimulator) sortOrderbooksBySymbol() []string {
	// Store all keys of the map into a slice
	sorted := make([]string, len(s.orderBooks))
	i := 0
	for symbol := range s.orderBooks {
		sorted[i] = symbol
		i++
	}
	sort.Strings(sorted)
	return sorted
}

func (s *binanceSimulator) sortAsksByPrice(symbol string) []float64 {
	sorted := make([]float64, len(s.orderBooks[symbol].asks))
	i := 0
	for price := range s.orderBooks[symbol].asks {
		sorted[i] = price
		i++
	}
	sort.Float64s(sorted)
	return sorted
}

func (s *binanceSimulator) sortBidsByPrice(symbol string) []float64 {
	sorted := make([]float64, len(s.orderBooks[symbol].asks))
	i := 0
	for price := range s.orderBooks[symbol].asks {
		sorted[i] = price
		i++
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(sorted)))
	return sorted

}

func (s *binanceSimulator) TakeSnapshot() (snapshot []Snapshot, err error) {
	snapshot = make([]Snapshot, 0, 10)
	// Take snapshots of subscribed channels
	sortedSubscribed := make([]string, len(s.subscribed))
	copy(sortedSubscribed, s.subscribed)
	sort.Strings(sortedSubscribed)
	for i, subChannel := range sortedSubscribed {
		subscribe := new(jsonstructs.BinanceSubscribe)
		subscribe.Initialize()
		// ID should not be 0
		subscribe.ID = i + 1
		subscribe.Params = []string{subChannel}

		subMarshaled, serr := json.Marshal(subscribe)
		if serr != nil {
			err = fmt.Errorf("subscribe marshal: %v", serr)
			return
		}
		snapshot = append(snapshot, Snapshot{
			Channel:  streamcommons.StateChannelSubscribed,
			Snapshot: subMarshaled,
		})
	}
	// Take snapshots of orderbooks
	for _, channel := range s.sortOrderbooksBySymbol() {
		symbol, _, serr := streamcommons.BinanceDecomposeChannel(channel)
		if serr != nil {
			err = serr
			return
		}
		memOrderbook := s.orderBooks[symbol]
		depth := new(jsonstructs.BinanceDepthREST)
		depth.Asks = make([][]string, len(memOrderbook.asks))
		for i, price := range s.sortAsksByPrice(symbol) {
			quantity := memOrderbook.asks[price]
			order := make([]string, 2)
			order[0] = strconv.FormatFloat(price, 'f', streamcommons.BinancePricePrecision, 64)
			order[1] = strconv.FormatFloat(quantity, 'f', streamcommons.BinanceQuantityPrecision, 64)
			depth.Asks[i] = order
		}
		depth.Bids = make([][]string, len(memOrderbook.bids))
		for i, price := range s.sortBidsByPrice(symbol) {
			quantity := memOrderbook.bids[price]
			order := make([]string, 2)
			order[0] = strconv.FormatFloat(price, 'f', streamcommons.BinancePricePrecision, 64)
			order[1] = strconv.FormatFloat(quantity, 'f', streamcommons.BinanceQuantityPrecision, 64)
			depth.Bids[i] = order
		}
		depth.LastUpdateID = memOrderbook.lastFinalUpdateID
		depthMarshaled, serr := json.Marshal(depth)
		if serr != nil {
			err = fmt.Errorf("orderbook marshal: %v", serr)
			return
		}
		snapshot = append(snapshot, Snapshot{
			Channel:  symbol + "@" + streamcommons.BinanceStreamRESTDepth,
			Snapshot: depthMarshaled,
		})
	}
	return
}

func newBinanceSimulator(channelFilter []string) *binanceSimulator {
	s := new(binanceSimulator)
	if channelFilter != nil {
		s.filterChannel = make(map[string]bool)
		for _, ch := range channelFilter {
			s.filterChannel[ch] = true
		}
	}
	s.idCh = make(map[int]string)
	s.orderBooks = make(map[string]*binanceOrderbook)
	return s
}
