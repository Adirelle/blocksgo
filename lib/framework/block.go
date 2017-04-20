package framework

import (
	"os"

	"github.com/Adirelle/blocksgo/lib/ipc"
)

// Block only implements the base block
type Block interface {
	Start() <-chan ipc.Block
	Stop()
}

// ClickableBlock handles click events
type ClickableBlock interface {
	Clicked(ipc.ClickEvent)
}

// SignalableBlock handles signals
type SignalableBlock interface {
	GetSignal() os.Signal
	Signaled(os.Signal)
}
