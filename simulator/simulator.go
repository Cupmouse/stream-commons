package simulator

import (
	"fmt"
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
	case "bitmex":
		return newBitmexSimulator(channels), nil
	case "bitflyer":
		return newBitflyerSimulator(channels), nil
	case "bitfinex":
		return newBitfinexSimulator(channels), nil
	case "binance":
		return newBinanceSimulator(channels), nil
	default:
		return nil, fmt.Errorf("snapshot for exchange %s is not supported", exchange)
	}
}
