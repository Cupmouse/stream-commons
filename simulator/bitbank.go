package simulator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/exchangedataset/streamcommons"
	"github.com/exchangedataset/streamcommons/jsonstructs"
)

type bitbankOrderbook struct {
	asks                map[float64]float64
	bids                map[float64]float64
	lastChangeTimestamp time.Time
}

type bitbankSimulator struct {
	filterChannel map[string]bool
	subscribed    []string
	orderbook     map[string]*bitbankOrderbook
}

func (s *bitbankSimulator) ProcessStart(line []byte) (err error) {
	return nil
}

func (s *bitbankSimulator) ProcessSend(line []byte) (channel string, err error) {
	channel = streamcommons.ChannelUnknown
	lineTrimmed := line[2:]
	sub := new(jsonstructs.BitbankSubscribe)
	serr := json.Unmarshal(lineTrimmed, sub)
	if serr != nil {
		err = fmt.Errorf("subscribe unmarshal: %v", serr)
		return
	}
	channel = sub[1]
	if s.filterChannel != nil {
		_, ok := s.filterChannel[channel]
		if !ok {
			return
		}
	}
	s.subscribed = append(s.subscribed, channel)
	if strings.HasPrefix(channel, "depth_whole_") {
		pair := channel[len("depth_whole_"):]
		orderbook := new(bitbankOrderbook)
		orderbook.asks = make(map[float64]float64)
		orderbook.bids = make(map[float64]float64)
		s.orderbook[pair] = orderbook
	} else if strings.HasPrefix(channel, "depth_diff_") {
		pair := channel[len("depth_diff_"):]
		orderbook := new(bitbankOrderbook)
		orderbook.asks = make(map[float64]float64)
		orderbook.bids = make(map[float64]float64)
		s.orderbook[pair] = orderbook
	}
	return
}

func (s *bitbankSimulator) processDepthSide(orderbook map[float64]float64, orders [][]string) error {
	for _, order := range orders {
		price, serr := strconv.ParseFloat(order[0], 64)
		if serr != nil {
			return fmt.Errorf("price: %v", serr)
		}
		amount, serr := strconv.ParseFloat(order[1], 64)
		if serr != nil {
			return fmt.Errorf("amount: %v", serr)
		}
		orderbook[price] = amount
	}
	return nil
}

func (s *bitbankSimulator) ProcessMessageWebSocket(line []byte) (channel string, err error) {
	channel = streamcommons.ChannelUnknown
	if line[0] == '0' {
		// Welcome message
		channel = "welcome"
		return
	}
	if bytes.HasPrefix(line, []byte("40")) {
		// Connected message?
		channel = "connected"
		return
	}
	// Raw line includes socketio constant 42
	lineTrimmed := line[2:]
	wrapper := new(jsonstructs.BitbankWrapper)
	serr := json.Unmarshal(lineTrimmed, wrapper)
	if serr != nil {
		err = fmt.Errorf("wrapper unmarshal: %v", serr)
		return
	}
	msgType := new(string)
	serr = json.Unmarshal(wrapper[0], msgType)
	if serr != nil {
		err = fmt.Errorf("msgType unmarshal: %v", serr)
		return
	}
	if *msgType != "message" {
		err = errors.New("wrapper but not message")
		return
	}
	root := new(jsonstructs.BitbankRoot)
	serr = json.Unmarshal(wrapper[1], root)
	if serr != nil {
		err = fmt.Errorf("root unmarshal: %v", serr)
		return
	}
	channel = root.RoomName
	if s.filterChannel != nil {
		_, ok := s.filterChannel[channel]
		if !ok {
			return
		}
	}
	if strings.HasPrefix(channel, "depth_whole_") {
		// Depth whole channel
		pair := channel[len("depth_whole_"):]
		depthWhole := new(jsonstructs.BitbankDepthWhole)
		serr := json.Unmarshal(root.Message, depthWhole)
		if serr != nil {
			err = fmt.Errorf("depth whole unmarshal: %v", serr)
			return
		}
		// Reset orderbook
		orderbook := new(bitbankOrderbook)
		orderbook.asks = make(map[float64]float64)
		orderbook.bids = make(map[float64]float64)
		err = s.processDepthSide(orderbook.asks, depthWhole.Asks)
		if err != nil {
			return
		}
		err = s.processDepthSide(orderbook.asks, depthWhole.Bids)
		if err != nil {
			return
		}
		orderbook.lastChangeTimestamp = unixMillisec(depthWhole.Timestamp)
		s.orderbook[pair] = orderbook
	} else if strings.HasPrefix(channel, "depth_diff_") {
		// Depth diff channel
		pair := channel[len("depth_diff_"):]
		depthDiff := new(jsonstructs.BitbankDepthDiff)
		serr := json.Unmarshal(root.Message, depthDiff)
		if serr != nil {
			err = fmt.Errorf("depth diff unmarshal: %v", serr)
			return
		}
		orderbook := s.orderbook[pair]
		err = s.processDepthSide(orderbook.asks, depthDiff.Asks)
		if err != nil {
			return
		}
		err = s.processDepthSide(orderbook.bids, depthDiff.Bids)
		if err != nil {
			return
		}
		orderbook.lastChangeTimestamp = unixMillisec(depthDiff.Timestamp)
	}
	return
}

func (s *bitbankSimulator) ProcessMessageChannelKnown(channel string, line []byte) (err error) {
	anoChannel, serr := s.ProcessMessageWebSocket(line)
	if anoChannel != channel {
		err = fmt.Errorf("channel differs")
	}
	if serr != nil {
		if err != nil {
			err = fmt.Errorf("%v, originally: %v", serr, err)
		} else {
			err = serr
		}
	}
	return
}

func (s *bitbankSimulator) ProcessState(channel string, line []byte) (err error) {
	if channel == streamcommons.StateChannelSubscribed {
		// Subscribed state
		subscribed := make([]string, 0)
		serr := json.Unmarshal(line, &subscribed)
		if serr != nil {
			return fmt.Errorf("state subscribed unmarshal: %v", serr)
		}
		if s.filterChannel == nil {
			s.subscribed = subscribed
		} else {
			for _, ch := range subscribed {
				_, ok := s.filterChannel[ch]
				if ok {
					s.subscribed = append(s.subscribed, ch)
				}
			}
		}
	}
	_, ok := s.filterChannel[channel]
	if !ok {
		return
	}
	if strings.HasPrefix(channel, "depth_whole_") {
		// Depth whole state (orderbook state)
		pair := channel[len("depth_whole_"):]
		orderbook := new(bitbankOrderbook)
		serr := json.Unmarshal(line, orderbook)
		if serr != nil {
			return fmt.Errorf("state depth whole unmarshal: %v", serr)
		}
		s.orderbook[pair] = orderbook
	}
	return
}

func (s *bitbankSimulator) TakeStateSnapshot() ([]Snapshot, error) {
	if s.filterChannel != nil {
		return nil, errors.New("channel filter is enabled")
	}
	snapshots := make([]Snapshot, 0, 100)
	// Subscribed channels snapshot
	subMar, serr := json.Marshal(s.subscribed)
	if serr != nil {
		return nil, fmt.Errorf("subscribed marshal: %v", serr)
	}
	snapshots = append(snapshots, Snapshot{
		Channel:  streamcommons.StateChannelSubscribed,
		Snapshot: subMar,
	})
	// Depth whole snapshots
	for pair, orderbook := range s.orderbook {
		orderbookMar, serr := json.Marshal(orderbook)
		if serr != nil {
			return nil, fmt.Errorf("orderbook marshal: %v", serr)
		}
		snapshots = append(snapshots, Snapshot{
			Channel:  "depth_whole_" + pair,
			Snapshot: orderbookMar,
		})
	}
	return snapshots, nil
}

func (s *bitbankSimulator) convertOrderbookSide(m map[float64]float64, reverse bool) [][]string {
	keys := make([]float64, len(m))
	i := 0
	for _, key := range m {
		keys[i] = key
		i++
	}
	if reverse {
		sort.Sort(sort.Reverse(sort.Float64Slice(keys)))
	} else {
		sort.Float64s(keys)
	}
	converted := make([][]string, len(m))
	for i, price := range keys {
		priceStr := strconv.FormatFloat(price, 'f', 8, 64)
		amountStr := strconv.FormatFloat(price, 'f', 8, 64)
		converted[i] = []string{priceStr, amountStr}
	}
	return converted
}

func (s *bitbankSimulator) TakeSnapshot() ([]Snapshot, error) {
	snapshots := make([]Snapshot, 0, 100)
	for pair, orderbook := range s.orderbook {
		depthWhole := new(jsonstructs.BitbankDepthWhole)
		depthWhole.Asks = s.convertOrderbookSide(orderbook.asks, false)
		depthWhole.Bids = s.convertOrderbookSide(orderbook.bids, true)
		depthWhole.Timestamp = orderbook.lastChangeTimestamp.UnixNano() / int64(time.Millisecond)
		depthWholeMar, serr := json.Marshal(depthWhole)
		if serr != nil {
			return nil, fmt.Errorf("depth whole marshal: %v", serr)
		}
		depthWholeLine := make([]byte, len(depthWholeMar)+2)
		depthWholeLine[0], depthWholeLine[1] = '4', '2'
		copy(depthWholeLine[2:], depthWholeMar)
		snapshots = append(snapshots, Snapshot{
			Channel:  "depth_whole_" + pair,
			Snapshot: depthWholeLine,
		})
	}
	return snapshots, nil
}

func newBitbankSimulator(filterChannel []string) *bitbankSimulator {
	s := new(bitbankSimulator)
	if filterChannel != nil {
		s.filterChannel = make(map[string]bool)
		for _, ch := range filterChannel {
			s.filterChannel[ch] = true
		}
	}
	s.orderbook = make(map[string]*bitbankOrderbook)
	s.subscribed = make([]string, 0)
	return s
}
