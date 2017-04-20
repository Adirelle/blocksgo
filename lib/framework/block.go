package framework

import "github.com/Adirelle/blocksgo/lib/ipc"

// Block only implements the base block
type Block interface {
	Start() <-chan ipc.Block
	Stop()
}

// BaseBlock can be used a base for concrete blocks
type BaseBlock struct {
	ipc.Block `yaml:",inline"`

	ch chan ipc.Block
}

// Start initializes the output channel
func (b *BaseBlock) Start() <-chan ipc.Block {
	b.ch = make(chan ipc.Block)
	return b.ch
}

// MakeDefaultBlock returns a preconfigured Block output.
func (b *BaseBlock) MakeDefaultBlock() ipc.Block {
	return b.Block
}

// Emit is used to send a Block to the output.
func (b *BaseBlock) Emit(bl ipc.Block) {
	b.ch <- bl
}

// Stop closes the output channel.
func (b *BaseBlock) Stop() {
	if b.ch != nil {
		close(b.ch)
		b.ch = nil
	}
}
