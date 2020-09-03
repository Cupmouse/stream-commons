package simulator

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

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
				if _, ok := s.filterChannel[subChannel]; ok {
					s.subscribed = append(s.subscribed, subChannel)
				}
			}
		}
		return
	}
	if s.filterChannel != nil {
		if _, ok := s.filterChannel[channel]; !ok {
			return
		}
	}
	return nil
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
	return &gen
}
