package framework

import (
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/Adirelle/blocksgo/lib/ipc"
)

type HandlerFunc func(value interface{})

type Application struct {
	Config

	line     ipc.StatusLine
	ticker   *time.Ticker
	cases    []reflect.SelectCase
	handlers []HandlerFunc
}

func (a *Application) Run(r io.Reader, w io.Writer) {
	a.start(r, w)

	a.loop()
}

func (a *Application) start(r io.Reader, w io.Writer) {
	fmt.Printf("%+v\n", a)
	a.ticker = time.NewTicker(a.Global.Interval)
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

	output, input := ipc.StartIPC(h, w, r)

	if h.ClickEvents {
		a.addRecvSelect(input, func(v interface{}) {
			a.dispatchClick(v.(ipc.ClickEvent))
		})
	}

	a.addRecvSelect(a.ticker.C, func(_ interface{}) {
		output <- a.line
	})
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

func (a *Application) loop() {
	for {
		i, v, ok := reflect.Select(a.cases)
		if !ok {
			break
		}
		a.handlers[i](v.Interface())
	}
}
