package framework

import (
	"io"
	"reflect"
	"time"

	"github.com/Adirelle/blocksgo/lib/ipc"
)

type HandlerFunc func(value interface{})

type Application struct {
	Config

	line     ipc.StatusLine
	cases    []reflect.SelectCase
	handlers []HandlerFunc
}

func (a *Application) Run(r io.Reader, w io.Writer) {
	output := a.start(r, w)

	a.loop(output)
}

func (a *Application) start(r io.Reader, w io.Writer) (output chan<- ipc.StatusLine) {
	a.line = make([]ipc.Block, len(a.Blocks))

	h := ipc.Header{Version: 1}

	for i, b := range a.Blocks {

		if _, handlesClicks := b.(ClickableBlock); handlesClicks {
			h.ClickEvents = true
		}

		ch := b.Start()
		a.addRecvSelect(ch, func(value interface{}) {
			a.line[i] = value.(ipc.Block)
		})
	}

	var input <-chan ipc.ClickEvent
	output, input = ipc.StartIPC(h, w, r)

	if h.ClickEvents {
		a.addRecvSelect(input, func(v interface{}) {
			a.dispatchClick(v.(ipc.ClickEvent))
		})
	}

	return
}

func (a *Application) addRecvSelect(ch interface{}, h HandlerFunc) {
	c := reflect.SelectCase{
		Chan: reflect.ValueOf(ch),
		Dir:  reflect.SelectRecv,
	}

	a.cases = append(a.cases, c)
	a.handlers = append(a.handlers, h)
}

func (a *Application) dispatchClick(e ipc.ClickEvent) {
	for _, b := range a.Blocks {
		if cb, handlesClicks := b.(ClickableBlock); handlesClicks && cb.AcceptClick(e) {
			cb.Clicked(e)
		}
	}
}

func (a *Application) loop(output chan<- ipc.StatusLine) {
	var lastUpdate time.Time

	for {
		if now := time.Now(); now.Sub(lastUpdate) >= a.Global.Interval {
			lastUpdate = now
			output <- a.line
		}

		i, v, ok := reflect.Select(a.cases)
		if !ok {
			break
		}
		a.handlers[i](v.Interface())
	}
}
