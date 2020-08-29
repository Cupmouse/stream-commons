package simulator

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/exchangedataset/streamcommons"
	"github.com/exchangedataset/streamcommons/jsonstructs"
)

type liquidSimulator struct {
	channelFilter map[string]bool
	subscribed    []string
}

func (s *liquidSimulator) ProcessStart(line []byte) (err error) {
	return nil
}

func (s *liquidSimulator) ProcessSend(line []byte) (channel string, err error) {
	channel = streamcommons.ChannelUnknown
	r := new(jsonstructs.LiquidMessageRoot)
	serr := json.Unmarshal(line, r)
	if serr != nil {
		err = fmt.Errorf("root unmarshal: %v", serr)
		return
	}
	d := new(jsonstructs.LiquidSubscribeData)
	serr = json.Unmarshal(r.Data, d)
	if serr != nil {
		err = fmt.Errorf("subscribe data unmarshal: %v", serr)
		return
	}
	channel = d.Channel
	return
}

func (s *liquidSimulator) ProcessMessageWebSocket(line []byte) (channel string, err error) {
	channel = streamcommons.ChannelUnknown
	r := new(jsonstructs.LiquidMessageRoot)
	serr := json.Unmarshal(line, r)
	if serr != nil {
		err = fmt.Errorf("root unmarshal: %v", serr)
		return
	}
	if r.Event == jsonstructs.LiquidConnectionEstablished {
		channel = streamcommons.LiquidChannelConnectionEstablished
		return
	}
	channel = *r.Channel
	if r.Event == jsonstructs.LiquidEventSubscriptionSucceeded {
		// Subscription successful
		if s.channelFilter != nil {
			_, ok := s.channelFilter[channel]
			if !ok {
				return
			}
		}
		s.subscribed = append(s.subscribed, channel)
	}
	return
}

func (s *liquidSimulator) ProcessMessageChannelKnown(channel string, line []byte) (err error) {
	anoChannel, serr := s.ProcessMessageWebSocket(line)
	if anoChannel != channel {
		err = errors.New("channel differs")
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

func (s *liquidSimulator) ProcessState(channel string, line []byte) (err error) {
	// Liquid has only subscribe state
	subscribed := make([]string, 0, 100)
	serr := json.Unmarshal(line, &subscribed)
	if serr != nil {
		return fmt.Errorf("subscribed unmarshal: %v", serr)
	}
	if s.channelFilter == nil {
		s.subscribed = subscribed
	} else {
		for _, ch := range subscribed {
			_, ok := s.channelFilter[ch]
			if ok {
				s.subscribed = append(s.subscribed, ch)
			}
		}
	}
	return nil
}

func (s *liquidSimulator) TakeStateSnapshot() ([]Snapshot, error) {
	if s.channelFilter != nil {
		return nil, errors.New("channel filter is enabled")
	}
	subm, serr := json.Marshal(s.subscribed)
	if serr != nil {
		return nil, fmt.Errorf("subscribed marshal: %v", serr)
	}
	return []Snapshot{
		Snapshot{
			Channel:  streamcommons.StateChannelSubscribed,
			Snapshot: subm,
		},
	}, nil
}

func (s *liquidSimulator) TakeSnapshot() ([]Snapshot, error) {
	snapshots := make([]Snapshot, len(s.subscribed))
	for i, ch := range s.subscribed {
		s := new(jsonstructs.LiquidSubscribeData)
		s.Channel = ch
		sm, serr := json.Marshal(s)
		if serr != nil {
			return nil, fmt.Errorf("subscribed marshal: %v", serr)
		}
		r := new(jsonstructs.LiquidMessageRoot)
		r.Channel = &ch
		r.Data = sm
		r.Event = jsonstructs.LiquidEventSubscriptionSucceeded
		rm, serr := json.Marshal(r)
		if serr != nil {
			return nil, fmt.Errorf("root marshal: %v", serr)
		}
		snapshots[i] = Snapshot{
			Channel:  ch,
			Snapshot: rm,
		}
	}
	return snapshots, nil
}

func newLiquidSimulator(channelFilter []string) *liquidSimulator {
	s := new(liquidSimulator)
	if channelFilter != nil {
		s.channelFilter = make(map[string]bool)
		for _, ch := range channelFilter {
			s.channelFilter[ch] = true
		}
	}
	s.subscribed = make([]string, 0, 100)
	return s
}
