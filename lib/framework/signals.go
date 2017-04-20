package framework

import "os"

// SignalableBlock handles signals
type SignalableBlock interface {
	GetSignal() os.Signal
	Signaled(os.Signal)
}
