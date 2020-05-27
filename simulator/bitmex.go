package simulator

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/exchangedataset/streamcommons/jsonstructs"
)

// BitmexChannelInfo is the channel name for info channel on Bitmex
const BitmexChannelInfo = "info"

// BitmexChannelError is the channel name for error channel on Bitmex
const BitmexChannelError = "error"

// BitmexSideSell is the string name for sell side
const BitmexSideSell = "Sell"

// BitmexSideBuy is the string name for buy side
const BitmexSideBuy = "Buy"

type bitmexOrderBookL2Element struct {
	price float64
	size  uint64
}

type bitmexSimulator struct {
	filterChannel map[string]bool
	subscribed    map[string]bool
	// map[symbol]map[side]map[id]
	orderBooks map[string]map[string]map[int64]bitmexOrderBookL2Element
}

func (g *bitmexSimulator) ProcessSend(line []byte) (channel string, err error) {
	// this should not be called
	return ChannelUnknown, nil
}

func (g *bitmexSimulator) processData(action string, dataSlice []jsonstructs.BitmexOrderBookL2DataElement) error {
	for _, data := range dataSlice {
		if action == "partial" || action == "insert" {
			sides, ok := g.orderBooks[data.Symbol]
			if !ok {
				// symbol is not yet pushed, create new map
				sides = make(map[string]map[int64]bitmexOrderBookL2Element)
				g.orderBooks[data.Symbol] = sides
			}
			ids, ok := sides[data.Side]
			if !ok {
				// map for this side is not prepared, create new one
				ids = make(map[int64]bitmexOrderBookL2Element)
				sides[data.Side] = ids
			}
			// set new element
			ids[data.ID] = bitmexOrderBookL2Element{price: data.Price, size: data.Size}

			// check logical error
			if data.Side == BitmexSideBuy {
				sellIDs, ok := sides[BitmexSideSell]
				// actually, it does not have to check if map exists, because range of nil map is no-op
				if ok {
					dueToRemove := make([]int64, 0, 5)
					for anoID, anoElem := range sellIDs {
						if anoElem.price < data.Price {
							// original order is buy, this order is sell but has lower price than original, weird
							dueToRemove = append(dueToRemove, anoID)
							fmt.Println("sell logical error:", data.Price, data.Size, anoElem.price, anoElem.size)
						}
					}
					for _, anoID := range dueToRemove {
						delete(sellIDs, anoID)
					}
				}
			} else {
				buyIDs, ok := sides[BitmexSideBuy]
				if ok {
					dueToRemove := make([]int64, 0, 5)
					for anoID, anoElem := range buyIDs {
						if anoElem.price > data.Price {
							// original order is sell, this order is buy but has higher price than original, weird
							dueToRemove = append(dueToRemove, anoID)
							fmt.Println("buy logical error:", data.Price, data.Size, anoElem.price, anoElem.size)
						}
					}
					for _, anoID := range dueToRemove {
						delete(buyIDs, anoID)
					}
				}
			}
		} else if action == "update" {
			// update for element, it can expect element to be there already,
			// so map is already prepared
			// map returns-by-value it needs to replace value after you updated it
			elem, ok := g.orderBooks[data.Symbol][data.Side][data.ID]
			if ok {
				elem.size = data.Size
				g.orderBooks[data.Symbol][data.Side][data.ID] = elem
			} else {
				fmt.Println("order id not found")
			}
		} else if action == "delete" {
			// delete element
			delete(g.orderBooks[data.Symbol][data.Side], data.ID)
		} else {
			return fmt.Errorf("unknown action type '%s'", action)
		}
	}
	return nil
}

func (g *bitmexSimulator) ProcessMessage(line []byte) (channel string, err error) {
	channel = ChannelUnknown

	// check if this message is a response to subscribe
	subscribe := jsonstructs.BitmexSubscribe{}
	err = json.Unmarshal(line, &subscribe)
	if err != nil {
		return
	}
	if subscribe.Success {
		// this is subscribe message
		channel = subscribe.Subscribe
		// check if this channel should be tracked
		if g.filterChannel == nil {
			g.subscribed[channel] = true
		} else {
			// filtering is enabled
			_, ok := g.filterChannel[channel]
			if ok {
				g.subscribed[channel] = true
			}
		}

		return
	}

	decoded := new(jsonstructs.BitmexRoot)
	err = json.Unmarshal(line, decoded)
	if err != nil {
		return
	}
	if decoded.Info != nil {
		channel = BitmexChannelInfo
		return
	}
	if decoded.Error != nil {
		channel = BitmexChannelError
		return
	}
	channel = decoded.Table

	if channel == "orderBookL2" {
		dataSlice := make([]jsonstructs.BitmexOrderBookL2DataElement, 0, 10)
		err = json.Unmarshal(decoded.Data, &dataSlice)
		if err != nil {
			return
		}

		err = g.processData(decoded.Action, dataSlice)

		return
	}
	// ignore other channels
	return
}

func (g *bitmexSimulator) ProcessState(channel string, line []byte) (err error) {
	if channel == StateChannelSubscribed {
		// add to subscribed
		subscribed := jsonstructs.BitmexStateSubscribed{}
		err = json.Unmarshal(line, &subscribed)
		if err != nil {
			return
		}
		for _, subscrCh := range subscribed {
			// record subscribed channel only if it is in target channel
			_, ok := g.filterChannel[subscrCh]
			if ok {
				g.subscribed[subscrCh] = true
			}
		}
		return
	}

	_, ok := g.filterChannel[channel]
	if !ok {
		return
	}

	if channel == "orderBookL2" {
		// process orderbook
		decoded := make([]jsonstructs.BitmexOrderBookL2DataElement, 0, 10)
		err = json.Unmarshal(line, &decoded)
		if err != nil {
			return
		}
		return g.processData("partial", decoded)
	}

	return
}

func sortBitmexSubscribe(m map[string]bool) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func sortBitmexOrderbooks(m map[string]map[string]map[int64]bitmexOrderBookL2Element) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func sortBitmexSides(m map[string]map[int64]bitmexOrderBookL2Element) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func sortBitmexID(m map[int64]bitmexOrderBookL2Element) []int64 {
	keys := make([]int64, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

func (g *bitmexSimulator) orderBookL2DataElements() []jsonstructs.BitmexOrderBookL2DataElement {
	// reconstruct raw-like json format bitmex sends
	data := make([]jsonstructs.BitmexOrderBookL2DataElement, 0, 10)
	for _, symbol := range sortBitmexOrderbooks(g.orderBooks) {
		sides := g.orderBooks[symbol]
		for _, side := range sortBitmexSides(sides) {
			ids := sides[side]
			for _, id := range sortBitmexID(ids) {
				elem := ids[id]
				data = append(data, jsonstructs.BitmexOrderBookL2DataElement{
					ID:     id,
					Price:  elem.price,
					Side:   side,
					Size:   elem.size,
					Symbol: symbol,
				})
			}
		}
	}
	return data
}

func (g *bitmexSimulator) TakeStateSnapshot() (snapshots []Snapshot, err error) {
	snapshots = make([]Snapshot, 0, 5)

	// list subscribed channels
	subList := make([]string, len(g.subscribed))
	for i, channel := range sortBitmexSubscribe(g.subscribed) {
		subList[i] = channel
	}
	var subListMarshaled []byte
	subListMarshaled, err = json.Marshal(subList)
	if err != nil {
		return nil, fmt.Errorf("error on json marshal: %s", err.Error())
	}
	snapshots = append(snapshots, Snapshot{Channel: StateChannelSubscribed, Snapshot: subListMarshaled})

	data := g.orderBookL2DataElements()
	var orderBookL2ElementsMarshaled []byte
	orderBookL2ElementsMarshaled, err = json.Marshal(data)
	if err != nil {
		return
	}
	snapshots = append(snapshots, Snapshot{Channel: "orderBookL2", Snapshot: orderBookL2ElementsMarshaled})

	return
}

func (g *bitmexSimulator) TakeSnapshot() (snapshots []Snapshot, err error) {
	snapshots = make([]Snapshot, 0, 5)

	// subscribe message
	for _, channel := range sortBitmexSubscribe(g.subscribed) {
		subscr := jsonstructs.BitmexSubscribe{}
		subscr.Initialize()
		subscr.Subscribe = channel

		var subscribeMarshaled []byte
		subscribeMarshaled, err = json.Marshal(subscr)
		if err != nil {
			return
		}

		snapshots = append(snapshots, Snapshot{Channel: channel, Snapshot: subscribeMarshaled})
	}

	_, ok := g.subscribed["orderBookL2"]
	if ok {
		// reconstruct raw-like json format bitmex sends
		data := g.orderBookL2DataElements()
		var dataMarshaled []byte
		dataMarshaled, err = json.Marshal(data)
		if err != nil {
			return
		}
		root := new(jsonstructs.BitmexRoot)
		root.Table = "orderBookL2"
		// partial means full orderbook snapshot
		root.Action = "partial"
		root.Data = json.RawMessage(dataMarshaled)
		var rootMarshaled []byte
		rootMarshaled, err = json.Marshal(root)
		if err != nil {
			return
		}
		snapshots = append(snapshots, Snapshot{Channel: "orderBookL2", Snapshot: rootMarshaled})
	}

	return
}

func newBitmexSimulator(filterChannels []string) Simulator {
	gen := bitmexSimulator{}
	if filterChannels != nil {
		gen.filterChannel = make(map[string]bool)
		for _, ch := range filterChannels {
			// this value will be ignored, filter will be applied to channel that has value in this map
			// value itself does not matter
			gen.filterChannel[ch] = true
		}
	}
	gen.subscribed = make(map[string]bool, 0)
	gen.orderBooks = make(map[string]map[string]map[int64]bitmexOrderBookL2Element)
	return &gen
}
