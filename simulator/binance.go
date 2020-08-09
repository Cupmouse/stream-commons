package simulator

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"

	"github.com/exchangedataset/streamcommons"
	"github.com/exchangedataset/streamcommons/jsonstructs"
)

type binanceOrderbook struct {
	asks map[float64]float64
	bids map[float64]float64
}

type binanceSimulator struct {
	filterChannel map[string]bool
	// IDs and its channel a client sent a subscription to
	idCh map[int]string
	// Slice of channels a client subscribed to and a server agreed
	subscribed []string
	// Last received FinalUpdateID to check for missing messages
	lastFinalUpdateID int64
	// Orderbook differences waiting to be applied for the arrival of a REST message
	// nil if a REST message has already been received
	differences map[string][]*jsonstructs.BinanceDepthStream
	// map[symbol]orderbook
	// Note: symbol is lower-cased one
	orderBooks map[string]*binanceOrderbook
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

func (s *binanceSimulator) processMessageDepth(channel string, depth *jsonstructs.BinanceDepthStream) (err error) {
	symbol, _, serr := streamcommons.BinanceDecomposeChannel(channel)
	if serr != nil {
		return serr
	}
	orderbook := s.orderBooks[symbol]
	if s.lastFinalUpdateID+1 != depth.FirstUpdateID {
		// There are missing messages that haven't been received
		err = fmt.Errorf("missing messages detected")
		return
	}
	for _, order := range depth.Asks {
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
			delete(orderbook.asks, price)
		} else {
			// Update order
			orderbook.asks[price] = quantity
		}
	}
	for _, order := range depth.Bids {
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
			delete(orderbook.bids, price)
		} else {
			orderbook.bids[price] = quantity
		}
	}
	return nil
}

func (s *binanceSimulator) ProcessMessageWebSocket(line []byte) (channel string, err error) {
	channel = streamcommons.ChannelUnknown
	subRes := new(jsonstructs.BinanceSubscribeResponse)
	err = json.Unmarshal(line, subRes)
	if err != nil {
		return
	}
	if subRes.ID != 0 {
		if subRes.Result != nil {
			err = fmt.Errorf("subscribe result != nil: %v", subRes.Result)
			return
		}
		// Subscribe message response
		channel = s.idCh[subRes.ID]
		if s.filterChannel != nil {
			_, ok := s.filterChannel[channel]
			if !ok {
				return
			}
		}
		// Add this channel to successfully subscribed channel list
		s.subscribed = append(s.subscribed, channel)
		symbol, stream, serr := streamcommons.BinanceDecomposeChannel(channel)
		if serr != nil {
			err = serr
			return
		}
		if stream == "depth@100ms" {
			// Create new orderbook in memory
			s.orderBooks[symbol] = new(binanceOrderbook)
			// Create a slice to store difference messages before a REST message arrives
			if _, ok := s.differences[symbol]; ok {
				err = errors.New("received subscribe confirmation twice")
				return
			}
			s.differences[symbol] = make([]*jsonstructs.BinanceDepthStream, 0, 1000)
		}
		return
	}
	// If its not a subscribe message, then must be in the root format
	root := new(jsonstructs.BinanceReponseRoot)
	channel = root.Stream
	symbol, stream, serr := streamcommons.BinanceDecomposeChannel(channel)
	if serr != nil {
		return
	}
	if stream == "depth100ms" {
		depth := new(jsonstructs.BinanceDepthStream)
		serr := json.Unmarshal(root.Data, depth)
		if serr != nil {
			err = fmt.Errorf("depth message unmarshal: %v", serr)
			return
		}
		differences := s.differences[symbol]
		if differences == nil {
			// Already received a REST message
			err = s.processMessageDepth(channel, depth)
			return
		}
		// Store this message into an slice
		s.differences[symbol] = append(differences, depth)
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
		// Apply orderbook differences previously received via WebSocket
		differences := s.differences[symbol]
		if differences == nil {
			err = errors.New("received REST twice or none")
			return
		}
		// Drop unneccesary messages
		i := 0
		for ; i < len(differences) && differences[i].FinalUpdateID <= depthRest.LastUpdateID; i++ {
		}
		if i == len(differences) {
			// No messages are stored
			err = fmt.Errorf("no messages stored")
			return
		}
		// First event should have this traits
		if differences[i].FirstUpdateID > depthRest.LastUpdateID+1 ||
			differences[i].FinalUpdateID <= depthRest.LastUpdateID+1 {
			return
		}
		for ; i < len(differences); i++ {
			// Apply all
			serr := s.processMessageDepth(symbol+"@depth@100ms", differences[i])
			if serr != nil {
				return fmt.Errorf("apply depth: %v", serr)
			}
		}
		s.differences[symbol] = nil
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
		depth.LastUpdateID = s.lastFinalUpdateID
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
	s.differences = make(map[string][]*jsonstructs.BinanceDepthStream)
	return s
}
