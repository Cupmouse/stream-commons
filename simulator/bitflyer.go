package simulator

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/exchangedataset/streamcommons/jsonstructs"
)

type bitflyerSimulator struct {
	filterChannel map[string]bool

	// id versus channel, channels subscribe message has been sent whether or not it actually subscribed map[messageID]channel
	idvch map[int]string
	// map[messageID]channel
	subscribed []string
	// map[channel]map[side]map[price]size
	orderBooks map[string]map[string]map[float64]float64
}

func (g *bitflyerSimulator) processOrders(channel string, message *jsonstructs.BitflyerBoardParamsMessage) {
	for _, ask := range message.Asks {
		if ask.Size == 0 {
			delete(g.orderBooks[channel]["asks"], ask.Price)
		} else {
			g.orderBooks[channel]["asks"][ask.Price] = ask.Size
		}
	}
	for _, bid := range message.Bids {
		if bid.Size == 0 {
			delete(g.orderBooks[channel]["bids"], bid.Price)
		} else {
			g.orderBooks[channel]["bids"][bid.Price] = bid.Size
		}
	}
}

func (g *bitflyerSimulator) ProcessSend(line []byte) (channel string, err error) {
	channel = ChannelUnknown

	subscribe := new(jsonstructs.BitflyerSubscribe)
	err = json.Unmarshal(line, subscribe)
	if err != nil {
		return
	}
	channel = subscribe.Params.Channel
	// store id and channel pair
	g.idvch[subscribe.ID] = channel

	return channel, nil
}

func (g *bitflyerSimulator) ProcessMessage(line []byte) (channel string, err error) {
	channel = ChannelUnknown

	subscribedUnmarshaled := new(jsonstructs.BitflyerSubscribed)
	err = json.Unmarshal(line, subscribedUnmarshaled)
	if err != nil {
		return
	}
	if subscribedUnmarshaled.Result {
		// response to a subscribe request
		channel = g.idvch[subscribedUnmarshaled.ID]
		if g.filterChannel != nil {
			_, ok := g.filterChannel[channel]
			if !ok {
				return
			}
		}
		g.subscribed = append(g.subscribed, channel)
		return
	}

	root := new(jsonstructs.BitflyerRoot)
	err = json.Unmarshal(line, root)
	if err != nil {
		return
	}
	channel = root.Params.Channel

	if strings.HasPrefix(channel, "lightning_board_") {
		if strings.HasPrefix(channel, "lightning_board_snapshot_") {
			// if channel is board snapshot, then remove all orders and start again
			g.orderBooks[channel] = make(map[string]map[float64]float64)
			g.orderBooks[channel]["asks"] = make(map[float64]float64)
			g.orderBooks[channel]["bids"] = make(map[float64]float64)
		}
		if g.orderBooks[channel] == nil {
			// snapshot is not received (lightning_board can be sent ealier than lightning_board_snapshot)
			return
		}
		message := new(jsonstructs.BitflyerBoardParamsMessage)
		err = json.Unmarshal(root.Params.Message, message)
		if err != nil {
			return
		}
		g.processOrders(channel, message)
	}
	return
}

func (g *bitflyerSimulator) ProcessState(channel string, line []byte) (err error) {
	if channel == StateChannelSubscribed {
		decoded := jsonstructs.BitflyerStateSubscribed{}
		err = json.Unmarshal(line, &decoded)
		if err != nil {
			return
		}
		for _, subscCh := range decoded {
			// add only target channel
			if g.filterChannel != nil {
				_, ok := g.filterChannel[subscCh]
				if ok {
					g.subscribed = append(g.subscribed, subscCh)
				}
			} else {
				g.subscribed = append(g.subscribed, subscCh)
			}
		}
		return
	}

	if g.filterChannel != nil {
		_, ok := g.filterChannel[channel]
		if !ok {
			return
		}
	}

	if strings.HasPrefix(channel, "lightning_board_snapshot_") {
		// dont forget to initalize orderbook map
		g.orderBooks[channel] = make(map[string]map[float64]float64)
		g.orderBooks[channel]["asks"] = make(map[float64]float64)
		g.orderBooks[channel]["bids"] = make(map[float64]float64)

		message := new(jsonstructs.BitflyerBoardParamsMessage)
		serr := json.Unmarshal(line, message)
		if serr != nil {
			return fmt.Errorf("unmarshal failed for params.message of state: %s", serr.Error())
		}
		g.processOrders(channel, message)
	}

	return nil
}

func sortBitflyerOrderbooks(m map[string]map[string]map[float64]float64) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func sortBitflyerPrice(m map[float64]float64) []float64 {
	keys := make([]float64, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Float64s(keys)
	return keys
}

func (g *bitflyerSimulator) takeOrderBookL2MessageSnapshot(channel string) (messageMarshaled []byte, err error) {
	memOrderBook := g.orderBooks[channel]
	message := new(jsonstructs.BitflyerBoardParamsMessage)

	memAsks := memOrderBook["asks"]
	memBids := memOrderBook["bids"]
	message.Asks = make([]jsonstructs.BitflyerBoardParamsMessageOrder, len(memAsks))
	i := 0
	for _, price := range sortBitflyerPrice(memAsks) {
		size := memAsks[price]
		message.Asks[i] = jsonstructs.BitflyerBoardParamsMessageOrder{Price: price, Size: size}
		i++
	}
	message.Bids = make([]jsonstructs.BitflyerBoardParamsMessageOrder, len(memBids))
	i = 0
	for _, price := range sortBitflyerPrice(memBids) {
		size := memBids[price]
		message.Bids[i] = jsonstructs.BitflyerBoardParamsMessageOrder{Price: price, Size: size}
		i++
	}
	messageMarshaled, err = json.Marshal(message)
	if err != nil {
		return
	}

	return
}

func (g *bitflyerSimulator) TakeStateSnapshot() (snapshots []Snapshot, err error) {
	snapshots = make([]Snapshot, 0, 5)

	// snapshot subscribed channels: list of channel names
	var subscribedMarshaled []byte
	subscribedMarshaled, err = json.Marshal(g.subscribed)
	if err != nil {
		return
	}
	snapshots = append(snapshots, Snapshot{Channel: StateChannelSubscribed, Snapshot: subscribedMarshaled})

	return
}

func (g *bitflyerSimulator) TakeSnapshot() (snapshots []Snapshot, err error) {
	snapshots = make([]Snapshot, 0, 5)

	// generate response message for subscribed channels
	for i, channel := range g.subscribed {
		subscr := new(jsonstructs.BitflyerSubscribed)
		subscr.Initialize()
		subscr.ID = i
		subscr.Result = true

		var marshaled []byte
		marshaled, err = json.Marshal(subscr)
		if err != nil {
			return
		}
		snapshots = append(snapshots, Snapshot{Channel: channel, Snapshot: marshaled})
	}

	for _, channel := range sortBitflyerOrderbooks(g.orderBooks) {
		book := new(jsonstructs.BitflyerRoot)
		// this is needed to initialize the constant value
		book.Initialize()
		book.Params.Channel = channel

		var messageMarshaled []byte
		messageMarshaled, err = g.takeOrderBookL2MessageSnapshot(channel)
		if err != nil {
			return
		}
		book.Params.Message = json.RawMessage(messageMarshaled)

		var bookMarshaled []byte
		bookMarshaled, err = json.Marshal(book)
		if err != nil {
			return
		}
		snapshots = append(snapshots, Snapshot{Channel: channel, Snapshot: bookMarshaled})
	}

	return
}

func newBitflyerSimulator(filterChannel []string) Simulator {
	gen := bitflyerSimulator{}
	if filterChannel != nil {
		gen.filterChannel = make(map[string]bool)
		for _, ch := range filterChannel {
			gen.filterChannel[ch] = true
		}
	}
	gen.idvch = make(map[int]string)
	gen.subscribed = make([]string, 0)
	gen.orderBooks = make(map[string]map[string]map[float64]float64)
	return &gen
}
