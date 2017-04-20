package framework

import "github.com/Adirelle/blocksgo/lib/ipc"

// Block only implements the base block
type Block interface {
	Start() <-chan ipc.Block
	Stop()
}

type BaseBlock struct {
	ipc.BlockIdentifier
	Color               string          `json:"color"`
	Background          string          `json:"background"`
	Border              string          `json:"border"`
	MinWidth            string          `json:"min_width"`
	Align               ipc.BlockAlign  `json:"align"`
	Separator           bool            `json:"separator"`
	SeparatorBlockWidth int             `json:"separator_block_width"`
	Markup              ipc.BlockMarkup `json:"markup,omitempty"`

	ch chan ipc.Block
}

func (b *BaseBlock) Start() <-chan ipc.Block {
	b.ch = make(chan ipc.Block)
	return b.ch
}

func (b *BaseBlock) MakeDefaultBlock() ipc.Block {
	return ipc.Block{
		BlockIdentifier:     b.BlockIdentifier,
		Color:               b.Color,
		Background:          b.Background,
		Border:              b.Border,
		MinWidth:            b.MinWidth,
		Align:               b.Align,
		Separator:           b.Separator,
		SeparatorBlockWidth: b.SeparatorBlockWidth,
		Markup:              b.Markup,
	}
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
