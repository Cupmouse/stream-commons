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
	// ProcessSend processes send message line and returns channel associate with it
	ProcessSend(line []byte) (channel string, err error)
	// MesssageLine processes message line returns channel associate with it
	ProcessMessage(line []byte) (channel string, err error)
	// StateLine processes state line, this should normally called first so it could reconstruct initial state
	// but not needed when it was used as read-time simulator (ProcessMessage will do that job) as all message from the beggining
	// will be given to this simulator.
	ProcessState(channel string, line []byte) (err error)
	// TakeStateSnapshot generates a snapshot as a state line, used only when it was used as a simulator for dumper
	TakeStateSnapshot() ([]Snapshot, error)
	// TakeSnapshot generates a snapshot, snapshots are in the format which is close to original exchange format
	TakeSnapshot() ([]Snapshot, error)
}

// GetSimulator return an appropriate generator for given parameters. if channels are not given, filtering will be disabled and all channel will be processed.
func GetSimulator(exchange string, channels []string) (Simulator, error) {
	switch exchange {
	case "bitmex":
		return newBitmexSimulator(channels), nil
	case "bitflyer":
		return newBitflyerSimulator(channels), nil
	case "bitfinex":
		return newBitfinexSimulator(channels), nil
	default:
		return nil, fmt.Errorf("snapshot for exchange %s is not supported", exchange)
	}
}
