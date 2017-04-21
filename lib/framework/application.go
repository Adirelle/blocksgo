package framework

import (
	"io"
	"time"

	"github.com/Adirelle/blocksgo/lib/ipc"
)

// Application is the base Struct
type Application struct {
	Config
	Select

	line   ipc.StatusLine
	input  <-chan ipc.ClickEvent
	output chan<- ipc.StatusLine
	ticker *time.Ticker
}

// Run runs the configured application, using the givne input and output streams.
func (a *Application) Run(r io.Reader, w io.Writer) {
	a.start(r, w)
	a.loop()
}

func (a *Application) start(r io.Reader, w io.Writer) {
	a.line = make([]ipc.Block, len(a.Blocks))

	h := ipc.Header{
		Version:     1,
		ClickEvents: hasClickableBlocks(a.Blocks),
	}
	a.output, a.input = ipc.StartIPC(h, w, r)

	if h.ClickEvents {
		a.OnRecv(a.input, func(v interface{}) {
			a.dispatchClick(v.(ipc.ClickEvent))
		})
	}

	for i, b := range a.Blocks {
		ch := b.Start()
		j := i
		a.OnRecv(ch, func(value interface{}) {
			a.updateBlock(j, value.(ipc.Block))
		})
	}

	a.ticker = time.NewTicker(a.Global.Interval)
	a.OnRecv(a.ticker.C, func(_ interface{}) {
		a.output <- a.line
	})

	return
}

func hasClickableBlocks(bs []Block) bool {
	for _, b := range bs {
		if _, handlesClicks := b.(ClickableBlock); handlesClicks {
			return true
		}
	}
	return false
}

func (a *Application) updateBlock(i int, b ipc.Block) {
	a.line[i] = b
}

func (a *Application) dispatchClick(e ipc.ClickEvent) {
	for _, b := range a.Blocks {
		if cb, handlesClicks := b.(ClickableBlock); handlesClicks && cb.AcceptClick(e) {
			cb.Clicked(e)
		}
	}
}

func (a *Application) loop() {
	for a.Select.Run() {
		// NOOP
	}
}
