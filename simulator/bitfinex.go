package simulator

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/exchangedataset/streamcommons"
	"github.com/exchangedataset/streamcommons/jsonstructs"
)

type bitfinexBookElement struct {
	count  uint64
	amount float64
}

// bitfinexSimulator generates a snapshot from data feeded
type bitfinexSimulator struct {
	filterChannel map[string]bool
	// map[chanID]channel
	idvch      map[int]string
	subscribed []int
	// map[channel]map[price]order
	orderBooks map[string]map[float64]bitfinexBookElement
}

func (s *bitfinexSimulator) ProcessStart(line []byte) error {
	return nil
}

// ProcessSend processes send message sent from client and returns associated channel
func (s *bitfinexSimulator) ProcessSend(line []byte) (channel string, err error) {
	channel = streamcommons.ChannelUnknown
	subscribe := new(jsonstructs.BitfinexSubscribe)
	err = json.Unmarshal(line, subscribe)
	if err != nil {
		return
	}
	channel = fmt.Sprintf("%s_%s", subscribe.Channel, subscribe.Symbol)
	return
}

func (s *bitfinexSimulator) processOrderBookL2Orders(channel string, ordersInterf interface{}) (err error) {
	memOrderBook, ok := s.orderBooks[channel]
	if !ok {
		s.orderBooks[channel] = make(map[float64]bitfinexBookElement)
		memOrderBook = s.orderBooks[channel]
	}

	var orders []jsonstructs.BitfinexBookOrder
	// orders will be flattened if there is only one order
	switch ordersInterf.(type) {
	case [][]interface{}:
		// probably made in golang
		orders = ordersInterf.([][]interface{})
	case []interface{}:
		ordersInterfs := ordersInterf.([]interface{})
		switch ordersInterfs[0].(type) {
		case float64:
			// only one order
			orders = []jsonstructs.BitfinexBookOrder{ordersInterfs}
		case []interface{}:
			orders = make([]jsonstructs.BitfinexBookOrder, len(ordersInterfs))
			for i, order := range ordersInterfs {
				orders[i] = order.(jsonstructs.BitfinexBookOrder)
			}
		default:
			return fmt.Errorf("invalid type for order: %s", reflect.TypeOf(ordersInterfs))
		}
	default:
		return fmt.Errorf("invalid type for ordersInterf: %s", reflect.TypeOf(ordersInterf))
	}
	for _, order := range orders {
		price := order[0].(float64)
		count := uint64(order[1].(float64))
		amount := order[2].(float64)
		if count == 0 {
			// delete order from orderbook
			delete(memOrderBook, price)
		} else {
			memOrderBook[price] = bitfinexBookElement{count: count, amount: amount}
			// removing logical error
			dueToRemove := make([]float64, 0, 5)
			for anoPrice, anoElem := range memOrderBook {
				if anoElem.amount*amount >= 0 {
					// this order is on the same side as original order
					continue
				}
				if (amount > 0 && anoPrice < price) || (amount < 0 && anoPrice > price) {
					// original order is buy and sell has lower price than the original, weird!
					// or sell and higher price
					dueToRemove = append(dueToRemove, anoPrice)
					fmt.Println("logical error:", price, amount, anoPrice, anoElem.amount)
				}
			}
			for _, anoPrice := range dueToRemove {
				delete(memOrderBook, anoPrice)
			}
		}
	}

	return nil
}

// ProcessMessageWebSocket processes message line from datasets and keep track of a internal state
func (s *bitfinexSimulator) ProcessMessageWebSocket(line []byte) (channel string, err error) {
	channel = streamcommons.ChannelUnknown
	subscribedStruct := new(jsonstructs.BitfinexSubscribed)
	// this might produce error as object could be an array
	// json.Unmarshal gives error when it tries to unmarshal array
	// into a struct
	err = json.Unmarshal(line, subscribedStruct)
	if err == nil {
		if subscribedStruct.Event != "subscribed" {
			// Event == error or info
			if subscribedStruct.Event == "error" {
				channel = fmt.Sprintf("%s_%s", subscribedStruct.Channel, subscribedStruct.Symbol)
			} else if subscribedStruct.Event == "info" {
				channel = "info"
			}
			return
		}
		// This is a subscribed response message from bitfinex
		// Dtore channel id and its name into map
		channel = fmt.Sprintf("%s_%s", subscribedStruct.Channel, subscribedStruct.Symbol)
		s.idvch[subscribedStruct.ChanID] = channel
		// Store to subscribed slice
		if s.filterChannel == nil {
			s.subscribed = append(s.subscribed, subscribedStruct.ChanID)
		} else {
			_, ok := s.filterChannel[channel]
			if ok {
				s.subscribed = append(s.subscribed, subscribedStruct.ChanID)
			}
		}
		return
	}
	decoded := make([]interface{}, 0, 5)
	err = json.Unmarshal(line, &decoded)
	if err != nil {
		return
	}
	chanID := int(decoded[0].(float64))
	channel = s.idvch[chanID]
	if strings.HasPrefix(channel, "book_") {
		switch decoded[1].(type) {
		case string:
			if decoded[1].(string) == "hb" {
				// this is heartbeat message, ignore
				return
			}
			return channel, fmt.Errorf("wrong string as a heartbeat: %s", decoded[1].(string))
		default:
			return channel, s.processOrderBookL2Orders(channel, decoded[1])
		}
	}
	// other channels are ignored as it does not have a state
	return
}

func (s *bitfinexSimulator) ProcessMessageChannelKnown(channel string, line []byte) error {
	wsChannel, serr := s.ProcessMessageWebSocket(line)
	if serr != nil {
		return serr
	}
	if wsChannel != channel {
		return fmt.Errorf("channel differs: %v, expected: %v", wsChannel, channel)
	}
	return nil
}

// ProcessState processes state line from a datasets
func (s *bitfinexSimulator) ProcessState(channel string, line []byte) (err error) {
	if channel == streamcommons.StateChannelSubscribed {
		decoded := make(jsonstructs.BitfinexStatusSubscribed)
		err = json.Unmarshal(line, &decoded)
		if err != nil {
			return
		}
		// register subscribed channels
		for ch, chanID := range decoded {
			s.idvch[chanID] = ch
			if s.filterChannel == nil {
				s.subscribed = append(s.subscribed, chanID)
			} else {
				_, ok := s.filterChannel[ch]
				if ok {
					s.subscribed = append(s.subscribed, chanID)
				}
			}
		}
		return
	}

	// from here we process message, if channel is not in filter-in map then return
	if s.filterChannel != nil {
		_, ok := s.filterChannel[channel]
		if !ok {
			return
		}
	}

	if strings.HasPrefix(channel, "book_") {
		// process book message
		// subscribed map have been filled before this
		decoded := make([]jsonstructs.BitfinexBookOrder, 0)
		err = json.Unmarshal(line, &decoded)
		if err != nil {
			return
		}
		err = s.processOrderBookL2Orders(channel, decoded)
		return
	}

	return
}

func sortBitfinexBooks(m map[string]map[float64]bitfinexBookElement) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func sortBitfinexBook(m map[float64]bitfinexBookElement) []float64 {
	keys := make([]float64, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Float64s(keys)
	return keys
}

// TakeStateSnapshot takes a snapshot of current state and return as state line
func (s *bitfinexSimulator) TakeStateSnapshot() (snapshots []Snapshot, err error) {
	if s.filterChannel != nil {
		// If channel filtering is enabled, this should not be called
		err = errors.New("channel filter is enabled")
		return
	}
	snapshots = make([]Snapshot, 0, 5)
	subscribed := make(jsonstructs.BitfinexStatusSubscribed)
	for _, chanID := range s.subscribed {
		subscribed[s.idvch[chanID]] = chanID
	}
	var subscribedMarshaled []byte
	subscribedMarshaled, err = json.Marshal(subscribed)
	if err != nil {
		return
	}
	snapshots = append(snapshots, Snapshot{
		Channel:  streamcommons.StateChannelSubscribed,
		Snapshot: subscribedMarshaled,
	})

	for _, channel := range sortBitfinexBooks(s.orderBooks) {
		memOrderBook := s.orderBooks[channel]

		orders := make([]jsonstructs.BitfinexBookOrder, len(memOrderBook))
		i := 0
		for _, price := range sortBitfinexBook(memOrderBook) {
			memOrder := memOrderBook[price]
			orders[i] = jsonstructs.BitfinexBookOrder{price, memOrder.count, memOrder.amount}
			i++
		}
		var ordersMarshaled []byte
		ordersMarshaled, err = json.Marshal(orders)
		if err != nil {
			return
		}
		snapshots = append(snapshots, Snapshot{Channel: channel, Snapshot: ordersMarshaled})
	}

	return
}

// TakeSnapshot takes a snapshot of current state and return
func (s *bitfinexSimulator) TakeSnapshot() (snapshots []Snapshot, err error) {
	snapshots = make([]Snapshot, 0, 5)

	// subscribed
	subSorted := make([]int, len(s.subscribed))
	copy(subSorted, s.subscribed)
	sort.Ints(subSorted)
	for _, chanID := range subSorted {
		channel := s.idvch[chanID]
		// this is needed to extract symbol and pair
		ind := strings.Index(channel, "_")
		bitfCh := channel[:ind]
		symbol := channel[ind+1:]

		subscribe := jsonstructs.BitfinexSubscribed{}
		// this will initialize event attribute
		subscribe.Initialize()
		subscribe.ChanID = chanID
		subscribe.Channel = bitfCh
		subscribe.Symbol = symbol
		subscribe.Pair = symbol[1:]

		var subscribeMarhsaled []byte
		subscribeMarhsaled, err = json.Marshal(subscribe)
		if err != nil {
			return
		}
		snapshots = append(snapshots, Snapshot{Channel: channel, Snapshot: subscribeMarhsaled})
	}

	// book channels
	for _, channel := range sortBitfinexBooks(s.orderBooks) {
		memOrderBook := s.orderBooks[channel]
		chanID := -1
		for subChanID, subChan := range s.idvch {
			if subChan == channel {
				chanID = subChanID
				break
			}
		}
		if chanID == -1 {
			err = fmt.Errorf("channel is not in idvch map: %s", channel)
			return
		}
		book := new(jsonstructs.BitfinexBook)

		book[0] = chanID
		// if there is only one order, flatten a slice of slice to just a slice
		// this is a official implementation on bitfinex, I personally hate it
		if len(memOrderBook) == 1 {
			for _, price := range sortBitfinexBook(memOrderBook) {
				memOrder := memOrderBook[price]
				// there is only one element in memOrderBook anyway
				book[1] = jsonstructs.BitfinexBookOrder{price, memOrder.count, memOrder.amount}
			}
		} else {
			orders := make([]jsonstructs.BitfinexBookOrder, len(memOrderBook))
			i := 0
			for _, price := range sortBitfinexBook(memOrderBook) {
				memOrder := memOrderBook[price]
				orders[i] = jsonstructs.BitfinexBookOrder{price, memOrder.count, memOrder.amount}
				i++
			}
			book[1] = orders
		}

		var bookMarshaled []byte
		bookMarshaled, err = json.Marshal(book)
		if err != nil {
			return
		}
		snapshots = append(snapshots, Snapshot{Channel: channel, Snapshot: bookMarshaled})
	}

	return
}

func newBitfinexSimulator(filterChannel []string) Simulator {
	gen := bitfinexSimulator{}
	if filterChannel != nil {
		gen.filterChannel = make(map[string]bool)
		for _, ch := range filterChannel {
			gen.filterChannel[ch] = true
		}
	}
	gen.idvch = make(map[int]string)
	gen.subscribed = make([]int, 0, 100)
	gen.orderBooks = make(map[string]map[float64]bitfinexBookElement)
	return &gen
}
