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

type binanceOrderbookState struct {
	Asks              [][2]float64                      `json:"asks"`
	Bids              [][2]float64                      `json:"bids"`
	IsLastSnapshot    bool                              `json:"isLastSnapshot"`
	LastFinalUpdateID int64                             `json:"lastFinalUpdateID"`
	Differences       []*jsonstructs.BinanceDepthStream `json:"differences"`
}

type binanceOrderbook struct {
	Asks map[float64]float64
	Bids map[float64]float64
	// True if immediate last operation was to construct the initial state
	IsLastSnapshot bool
	// Last received FinalUpdateID to check for missing messages
	LastFinalUpdateID int64
	// If this orderbook already received an initial state
	ReceivedREST bool
	// Orderbook differences waiting to be applied for the arrival of a REST message
	// nil if a REST message has already been received
	Differences []*jsonstructs.BinanceDepthStream
}

type binanceSimulator struct {
	channelFilter map[string]bool
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
		return fmt.Errorf("ProcessStart: %v", serr)
	}
	query := u.Query()
	channels := strings.Split(query.Get("streams"), "/")
	for _, ch := range channels {
		// Apply filter
		if s.channelFilter != nil {
			if _, ok := s.channelFilter[ch]; !ok {
				continue
			}
		}
		// Add this channel to successfully subscribed channel list
		s.subscribed = append(s.subscribed, ch)
		symbol, stream, serr := streamcommons.BinanceDecomposeChannel(ch)
		if serr != nil {
			return fmt.Errorf("ProcessStart: %v", serr)
		}
		if stream == streamcommons.BinanceStreamDepth {
			if _, ok := s.orderBooks[symbol]; ok {
				return errors.New("ProcessStart: received subscribe confirmation twice")
			}
			// Create new orderbook in memory
			orderbook := new(binanceOrderbook)
			orderbook.Asks = make(map[float64]float64, 10000)
			orderbook.Bids = make(map[float64]float64, 10000)
			// Create a slice to store difference messages before a REST message arrives
			orderbook.Differences = make([]*jsonstructs.BinanceDepthStream, 0, 1000)
			s.orderBooks[symbol] = orderbook
			if serr != nil {
				return fmt.Errorf("ProcessStart: %v", serr)
			}
		}
	}
	return nil
}

func (s *binanceSimulator) ProcessSend(line []byte) (channel string, err error) {
	// defer func() {
	// 	if err != nil {
	// 		err = fmt.Errorf("ProcessSend: %v", err)
	// 	}
	// }()
	// channel = streamcommons.ChannelUnknown
	// sub := new(jsonstructs.BinanceSubscribe)
	// serr := json.Unmarshal(line, sub)
	// if serr != nil {
	// 	err = fmt.Errorf("subscribe unmarshal: %v", serr)
	// 	return
	// }
	// if len(sub.Params) != 1 {
	// 	// Subscription have to be done one at the time
	// 	err = errors.New("len(subscribe.Params) != 1")
	// 	return
	// }
	// if sub.ID == 0 {
	// 	err = errors.New("use of 0 as subscription id")
	// }
	// channel = sub.Params[0]
	// s.idCh[sub.ID] = channel
	return streamcommons.ChannelUnknown, nil
}

func binanceProcessSide(asks [][]string, m map[float64]float64) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("binanceProcessSide: %v", err)
		}
	}()
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

func (s *binanceSimulator) processMessageDepth(symbol string, depth *jsonstructs.BinanceDepthStream) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("processMessageDepth: %v", err)
		}
	}()
	orderbook := s.orderBooks[symbol]
	if orderbook.IsLastSnapshot {
		// First event should have this traits
		if depth.FirstUpdateID > orderbook.LastFinalUpdateID+1 ||
			depth.FinalUpdateID < orderbook.LastFinalUpdateID+1 {
			err = fmt.Errorf("first difference's updateID out of range")
			return
		}
		orderbook.IsLastSnapshot = false
	} else {
		if orderbook.LastFinalUpdateID+1 != depth.FirstUpdateID {
			// There are missing messages that haven't been received
			err = fmt.Errorf("missing messages detected")
			return
		}
	}
	err = binanceProcessSide(depth.Asks, orderbook.Asks)
	if err != nil {
		return
	}
	err = binanceProcessSide(depth.Bids, orderbook.Bids)
	if err != nil {
		return
	}
	orderbook.LastFinalUpdateID = depth.FinalUpdateID
	return nil
}

func (s *binanceSimulator) ProcessMessageWebSocket(line []byte) (channel string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("ProcessMessageWebSocket: %v", err)
		}
	}()
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
	if stream == streamcommons.BinanceStreamDepth {
		_, ok := s.channelFilter[symbol+"@"+streamcommons.BinanceStreamRESTDepth]
		if !ok {
			// Don't have to be tracked
			return
		}
		depth := new(jsonstructs.BinanceDepthStream)
		serr := json.Unmarshal(root.Data, depth)
		if serr != nil {
			err = fmt.Errorf("depth message unmarshal: %v", serr)
			return
		}
		orderbook := s.orderBooks[symbol]
		if orderbook.ReceivedREST {
			// Already received a REST message
			err = s.processMessageDepth(symbol, depth)
			return
		}
		if len(orderbook.Differences) > 100 {
			err = fmt.Errorf("too much stored difference: %v", symbol)
			return
		}
		// Store this message into an slice
		orderbook.Differences = append(orderbook.Differences, depth)
	}
	// Ignore other channels
	return
}

func (s *binanceSimulator) ProcessMessageChannelKnown(channel string, line []byte) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("ProcessMessageChannelKnown: %v", err)
		}
	}()
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
		if orderbook.LastFinalUpdateID != 0 {
			err = errors.New("received REST twice")
			return
		}
		orderbook.IsLastSnapshot = true
		orderbook.LastFinalUpdateID = depthRest.LastUpdateID
		err = binanceProcessSide(depthRest.Asks, orderbook.Asks)
		if err != nil {
			return
		}
		err = binanceProcessSide(depthRest.Bids, orderbook.Bids)
		if err != nil {
			return
		}
		// Apply orderbook differences previously received via WebSocket
		differences := orderbook.Differences
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
			serr := s.processMessageDepth(symbol, differences[i])
			if serr != nil {
				return fmt.Errorf("apply depth: %v", serr)
			}
		}
		// To free memory space
		orderbook.Differences = nil
		orderbook.ReceivedREST = true
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
	defer func() {
		if err != nil {
			err = fmt.Errorf("ProcessState: %v", err)
		}
	}()
	if channel == streamcommons.StateChannelSubscribed {
		subscribed := make([]string, 0, 100)
		// Read subscribed state
		serr := json.Unmarshal(line, &subscribed)
		if serr != nil {
			err = fmt.Errorf("subscribed state unmarshal: %v", serr)
			return
		}
		if s.channelFilter == nil {
			// Filter is disabled, add all
			s.subscribed = subscribed
		} else {
			// Apply filter to it
			for _, stateChannel := range subscribed {
				if _, ok := s.channelFilter[stateChannel]; ok {
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
	if stream == streamcommons.BinanceStreamRESTDepth {
		// State for rest depth
		if s.channelFilter != nil {
			// Apply filter
			if _, ok := s.channelFilter[channel]; !ok {
				return
			}
		}
		state := new(binanceOrderbookState)
		serr := json.Unmarshal(line, state)
		if serr != nil {
			err = fmt.Errorf("orderbook unmarshal: %v", serr)
			return
		}
		ob := new(binanceOrderbook)
		ob.Asks = make(map[float64]float64)
		for _, arr := range state.Asks {
			ob.Asks[arr[0]] = arr[1]
		}
		ob.Bids = make(map[float64]float64)
		for _, arr := range state.Bids {
			ob.Bids[arr[0]] = arr[1]
		}
		ob.Differences = state.Differences
		ob.ReceivedREST = state.LastFinalUpdateID != 0
		ob.IsLastSnapshot = state.IsLastSnapshot
		ob.LastFinalUpdateID = state.LastFinalUpdateID
		s.orderBooks[symbol] = ob
		return
	}
	err = fmt.Errorf("unknown stream name: %v", stream)
	return
}

func (s *binanceSimulator) TakeStateSnapshot() (snapshots []Snapshot, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("TakeStateSnapshot: %v", err)
		}
	}()
	if s.channelFilter != nil {
		// If channel filtering is enabled, this should not be called
		err = errors.New("channel filter is enabled")
		return
	}
	snapshots = make([]Snapshot, 0, 100)
	// Take a snapshot of subscribed channel list
	subMarshaled, serr := json.Marshal(s.subscribed)
	if serr != nil {
		err = fmt.Errorf("subscribed marshal: %v", serr)
		return
	}
	snapshots = append(snapshots, Snapshot{
		Channel:  streamcommons.StateChannelSubscribed,
		Snapshot: subMarshaled,
	})
	// Take snapshots of orderbooks
	for symbol, orderbook := range s.orderBooks {
		state := new(binanceOrderbookState)
		state.Asks = make([][2]float64, len(orderbook.Asks))
		i := 0
		for price, quantity := range orderbook.Asks {
			state.Asks[i] = [2]float64{price, quantity}
			i++
		}
		state.Bids = make([][2]float64, len(orderbook.Bids))
		i = 0
		for price, quantity := range orderbook.Bids {
			state.Bids[i] = [2]float64{price, quantity}
			i++
		}
		state.Differences = orderbook.Differences
		state.IsLastSnapshot = orderbook.IsLastSnapshot
		state.LastFinalUpdateID = orderbook.LastFinalUpdateID
		sm, serr := json.Marshal(state)
		if serr != nil {
			err = fmt.Errorf("orderbook marshal: %v", serr)
			return
		}
		snapshots = append(snapshots, Snapshot{
			Channel:  symbol + "@" + streamcommons.BinanceStreamRESTDepth,
			Snapshot: sm,
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
	sorted := make([]float64, len(s.orderBooks[symbol].Asks))
	i := 0
	for price := range s.orderBooks[symbol].Asks {
		sorted[i] = price
		i++
	}
	sort.Float64s(sorted)
	return sorted
}

func (s *binanceSimulator) sortBidsByPrice(symbol string) []float64 {
	sorted := make([]float64, len(s.orderBooks[symbol].Bids))
	i := 0
	for price := range s.orderBooks[symbol].Bids {
		sorted[i] = price
		i++
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(sorted)))
	return sorted
}

func (s *binanceSimulator) TakeSnapshot() (snapshot []Snapshot, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("TakeSnapshot: %v", err)
		}
	}()
	snapshot = make([]Snapshot, 0, 10)
	// Take snapshots of subscribed channels
	sortedSubscribed := make([]string, len(s.subscribed))
	copy(sortedSubscribed, s.subscribed)
	sort.Strings(sortedSubscribed)
	for i, subChannel := range sortedSubscribed {
		subResp := new(jsonstructs.BinanceSubscribeResponse)
		// ID should not be 0
		subResp.ID = i + 1
		subMarshaled, serr := json.Marshal(subResp)
		if serr != nil {
			err = fmt.Errorf("subscribed marshal: %v", serr)
			return
		}
		snapshot = append(snapshot, Snapshot{
			Channel:  subChannel,
			Snapshot: subMarshaled,
		})
	}
	// Take snapshots of orderbooks
	for _, symbol := range s.sortOrderbooksBySymbol() {
		_, ok := s.channelFilter[symbol+"@"+streamcommons.BinanceStreamRESTDepth]
		if !ok {
			continue
		}
		memOrderbook := s.orderBooks[symbol]
		depth := new(jsonstructs.BinanceDepthREST)
		depth.Asks = make([][]string, len(memOrderbook.Asks))
		for i, price := range s.sortAsksByPrice(symbol) {
			quantity := memOrderbook.Asks[price]
			order := make([]string, 2)
			order[0] = strconv.FormatFloat(price, 'f', streamcommons.BinancePricePrecision, 64)
			order[1] = strconv.FormatFloat(quantity, 'f', streamcommons.BinanceQuantityPrecision, 64)
			depth.Asks[i] = order
		}
		depth.Bids = make([][]string, len(memOrderbook.Bids))
		for i, price := range s.sortBidsByPrice(symbol) {
			quantity := memOrderbook.Bids[price]
			order := make([]string, 2)
			order[0] = strconv.FormatFloat(price, 'f', streamcommons.BinancePricePrecision, 64)
			order[1] = strconv.FormatFloat(quantity, 'f', streamcommons.BinanceQuantityPrecision, 64)
			depth.Bids[i] = order
		}
		depth.LastUpdateID = memOrderbook.LastFinalUpdateID
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
		s.channelFilter = make(map[string]bool)
		for _, ch := range channelFilter {
			s.channelFilter[ch] = true
		}
	}
	s.idCh = make(map[int]string)
	s.orderBooks = make(map[string]*binanceOrderbook)
	return s
}
