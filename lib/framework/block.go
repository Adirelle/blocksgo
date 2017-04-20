package framework

import "github.com/Adirelle/blocksgo/lib/ipc"

// Block only implements the base block
type Block interface {
	Start() <-chan ipc.Block
	Stop()
}

type BaseBlock struct {
	ipc.Block `yaml:",inline"`

	ch chan ipc.Block
}

func (b *BaseBlock) Start() <-chan ipc.Block {
	b.ch = make(chan ipc.Block)
	return b.ch
}

func (b *BaseBlock) MakeDefaultBlock() ipc.Block {
	return b.Block
}

func (b *BaseBlock) Emit(bl ipc.Block) {
	b.ch <- bl
}

func (b *BaseBlock) Stop() {
	if b.ch != nil {
		close(b.ch)
		b.ch = nil
	}
}
