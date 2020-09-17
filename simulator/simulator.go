package simulator

import (
	"fmt"
	"strings"

	"github.com/exchangedataset/streamcommons"
)

// Snapshot is individual snapshot for a channel.
type Snapshot struct {
	Channel  string
	Snapshot []byte
}

// Simulator takes in message from exchange and chases the state and generates snapshot.
type Simulator interface {
	// ProcessStart processes start line who usually include connection URL.
	ProcessStart(line []byte) (err error)
	// ProcessSend processes send message line and returns the channel associated with.
	ProcessSend(line []byte) (channel string, err error)
	// ProcessMessageWebSocket processes a message from WebSocket and returns the channel associated with.
	// Other messages should not be inputted.
	ProcessMessageWebSocket(line []byte) (channel string, err error)
	// ProcessMessageChannelKnown is the same as `ProcessMessageWebSocket`, but should always be used
	// instead if the caller knows the channel from which the line is produced.
	// As a message, both message via REST or WebSocket is expected.
	ProcessMessageChannelKnown(channel string, line []byte) (err error)
	// ProcessState processes a state line, this should normally be called first
	// when it is used for dump file to reconstruct the initial state.
	// This should not be called when used for a real-time dump.
	ProcessState(channel string, line []byte) (err error)
	// TakeStateSnapshot generates a state line, used only when the simulator
	// is being used as the simulator for dump.
	TakeStateSnapshot() ([]Snapshot, error)
	// TakeSnapshot generates a snapshot, snapshots are in the format which is
	// close to original exchange format.
	TakeSnapshot() ([]Snapshot, error)
}

// GetSimulator return an appropriate generator for the given parameters. if channels are not given,
// filtering will be disabled subsequently all channel will be processed.
func GetSimulator(exchange string, channels []string) (Simulator, error) {
	switch exchange {
	case streamcommons.ExchangeBitmex:
		return newBitmexSimulator(channels), nil
	case streamcommons.ExchangeBitflyer:
		return newBitflyerSimulator(channels), nil
	case streamcommons.ExchangeBitfinex:
		return newBitfinexSimulator(channels), nil
	case streamcommons.ExchangeBinance:
		return newBinanceSimulator(channels), nil
	case streamcommons.ExchangeLiquid:
		return newLiquidSimulator(channels), nil
	default:
		return nil, fmt.Errorf("snapshot for exchange %s is not supported", exchange)
	}
}

// ToSimulatorChannel converts raw channels (user specified) to simulator channels.
func ToSimulatorChannel(exchange string, rawChannels []string) []string {
	switch exchange {
	case streamcommons.ExchangeBitmex:
		// Construct unique list of raw channels
		set := make(map[string]bool)
		for _, ch := range rawChannels {
			// Take prefix from full channel name
			// Eg. orderBookL2_XBTUSD -> orderBookL2
			if ri := strings.IndexRune(ch, '_'); ri != -1 {
				set[ch[:ri]] = true
			} else {
				set[ch] = true
			}
		}
		list := make([]string, len(set))
		i := 0
		for ch := range set {
			list[i] = ch
			i++
		}
		return list
	default:
		return rawChannels
	}
}
