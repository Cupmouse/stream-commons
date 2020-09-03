package simulator

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/exchangedataset/streamcommons"
	"github.com/exchangedataset/streamcommons/jsonstructs"
)

type bitflyerSimulator struct {
	filterChannel map[string]bool

	// id versus channel, channels subscribe message has been sent whether or not it actually subscribed map[messageID]channel
	// FIXME idvch is not preserved: should add new state line
	idvch map[int]string
	// map[messageID]channel
	subscribed []string
	// map[channel]map[side]map[price]size
	orderBooks map[string]map[string]map[float64]float64
}

func (s *bitflyerSimulator) ProcessStart(line []byte) error {
	return nil
}

func (s *bitflyerSimulator) processOrders(channel string, message *jsonstructs.BitflyerBoardParamsMessage) {
	for _, ask := range message.Asks {
		if ask.Size == 0 {
			delete(s.orderBooks[channel]["asks"], ask.Price)
		} else {
			s.orderBooks[channel]["asks"][ask.Price] = ask.Size
		}
	}
	for _, bid := range message.Bids {
		if bid.Size == 0 {
			delete(s.orderBooks[channel]["bids"], bid.Price)
		} else {
			s.orderBooks[channel]["bids"][bid.Price] = bid.Size
		}
	}
}

func (s *bitflyerSimulator) ProcessSend(line []byte) (channel string, err error) {
	channel = streamcommons.ChannelUnknown

	subscribe := new(jsonstructs.BitflyerSubscribe)
	err = json.Unmarshal(line, subscribe)
	if err != nil {
		return
	}
	channel = subscribe.Params.Channel
	// store id and channel pair
	s.idvch[subscribe.ID] = channel

	return channel, nil
}

func (s *bitflyerSimulator) ProcessMessageWebSocket(line []byte) (channel string, err error) {
	channel = streamcommons.ChannelUnknown

	subscribedUnmarshaled := new(jsonstructs.BitflyerSubscribed)
	err = json.Unmarshal(line, subscribedUnmarshaled)
	if err != nil {
		return
	}
	if subscribedUnmarshaled.Result {
		// response to a subscribe request
		channel = s.idvch[subscribedUnmarshaled.ID]
		if s.filterChannel != nil {
			_, ok := s.filterChannel[channel]
			if !ok {
				return
			}
		}
		s.subscribed = append(s.subscribed, channel)
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
			s.orderBooks[channel] = make(map[string]map[float64]float64)
			s.orderBooks[channel]["asks"] = make(map[float64]float64)
			s.orderBooks[channel]["bids"] = make(map[float64]float64)
		}
		if s.orderBooks[channel] == nil {
			// snapshot is not received (lightning_board can be sent ealier than lightning_board_snapshot)
			return
		}
		message := new(jsonstructs.BitflyerBoardParamsMessage)
		err = json.Unmarshal(root.Params.Message, message)
		if err != nil {
			return
		}
		s.processOrders(channel, message)
	}
	return
}

func (s *bitflyerSimulator) ProcessMessageChannelKnown(channel string, line []byte) error {
	wsChannel, serr := s.ProcessMessageWebSocket(line)
	if serr != nil {
		return serr
	}
	if wsChannel != channel {
		return fmt.Errorf("channel differs: %v, expected: %v", wsChannel, channel)
	}
	return nil
}

func (s *bitflyerSimulator) ProcessState(channel string, line []byte) (err error) {
	if channel == streamcommons.StateChannelSubscribed {
		subscribed := make(jsonstructs.BitflyerStateSubscribed, 0, 50)
		err = json.Unmarshal(line, &subscribed)
		if err != nil {
			return
		}
		if s.filterChannel == nil {
			s.subscribed = subscribed
		} else {
			// Add only target channels
			for _, subChannel := range subscribed {
				_, ok := s.filterChannel[subChannel]
				if ok {
					s.subscribed = append(s.subscribed, subChannel)
				}
			}
		}
		return
	}

	if s.filterChannel != nil {
		_, ok := s.filterChannel[channel]
		if !ok {
			return
		}
	}

	if strings.HasPrefix(channel, "lightning_board_snapshot_") {
		// dont forget to initalize orderbook map
		s.orderBooks[channel] = make(map[string]map[float64]float64)
		s.orderBooks[channel]["asks"] = make(map[float64]float64)
		s.orderBooks[channel]["bids"] = make(map[float64]float64)

		message := new(jsonstructs.BitflyerBoardParamsMessage)
		serr := json.Unmarshal(line, message)
		if serr != nil {
			return fmt.Errorf("unmarshal failed for params.message of state: %s", serr.Error())
		}
		s.processOrders(channel, message)
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

func (s *bitflyerSimulator) takeOrderBookL2MessageSnapshot(channel string) (messageMarshaled []byte, err error) {
	memOrderBook := s.orderBooks[channel]
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

func (s *bitflyerSimulator) TakeStateSnapshot() (snapshots []Snapshot, err error) {
	if s.filterChannel != nil {
		// If channel filtering is enabled, this should not be called
		err = errors.New("channel filter is enabled")
		return
	}
	snapshots = make([]Snapshot, 0, 5)

	// snapshot subscribed channels: list of channel names
	var subscribedMarshaled []byte
	subscribedMarshaled, err = json.Marshal(s.subscribed)
	if err != nil {
		return
	}
	snapshots = append(snapshots, Snapshot{
		Channel:  streamcommons.StateChannelSubscribed,
		Snapshot: subscribedMarshaled,
	})

	return
}

func (s *bitflyerSimulator) TakeSnapshot() (snapshots []Snapshot, err error) {
	snapshots = make([]Snapshot, 0, 5)

	// generate response message for subscribed channels
	sortedSubscibed := make([]string, len(s.subscribed))
	copy(sortedSubscibed, s.subscribed)
	sort.Strings(sortedSubscibed)
	for i, channel := range sortedSubscibed {
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

	for _, channel := range sortBitflyerOrderbooks(s.orderBooks) {
		_, ok := s.filterChannel[channel]
		if !ok {
			// Filter out
			continue
		}
		book := new(jsonstructs.BitflyerRoot)
		// This is needed to initialize the constant value
		book.Initialize()
		book.Params.Channel = channel

		var messageMarshaled []byte
		messageMarshaled, err = s.takeOrderBookL2MessageSnapshot(channel)
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
