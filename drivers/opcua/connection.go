package opcua

import (
	"sync"
	"sync/atomic"
)

type State int32

const (
	StateDisconnected State = iota
	StateConnecting
	StateConnected
)

func (s State) String() string {
	switch s {
	case StateDisconnected:
		return "Disconnected"
	case StateConnecting:
		return "Connecting"
	case StateConnected:
		return "Connected"
	default:
		return "Unknown"
	}
}

type BaseAdapter struct {
	state int32
	mu    sync.RWMutex
}

func (a *BaseAdapter) State() State {
	return State(atomic.LoadInt32(&a.state))
}

func (a *BaseAdapter) SetState(s State) {
	atomic.StoreInt32(&a.state, int32(s))
}

func (a *BaseAdapter) IsConnected() bool {
	return a.State() == StateConnected
}