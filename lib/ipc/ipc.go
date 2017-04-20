package ipc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// The Header is sent on initialization
type Header struct {
	Version     int  `json:"version"`
	StopSignal  int  `json:"stop_signal,omitempty"`
	ContSignal  int  `json:"cont_signal,omitempty"`
	ClickEvents bool `json:"click_events,omitempty"`
}

// BlockMarkup describes the type of markup used in block text
type BlockMarkup string

// BlockAlign describes how the align the text content if it is
// narrowed than the minimal width
type BlockAlign string

const (
	NoMarkup    BlockMarkup = "none"
	PangoMarkup BlockMarkup = "pango"

	AlignLeft   BlockAlign = "left"
	AlignCenter BlockAlign = "center"
	AlignRight  BlockAlign = "right"
)

// BlockIdentifier uniquely identify a Block
type BlockIdentifier struct {
	Name     string `json:"name,omitempty"`
	Instance string `json:"instance,omitempty"`
}

// A Block describes a part of the i3bar
type Block struct {
	BlockIdentifier     `yaml:",inline"`
	FullText            string      `json:"full_text"`
	ShortText           string      `json:"short_text,omitempty"`
	Color               string      `json:"color,omitempty"`
	Background          string      `json:"background,omitempty"`
	Border              string      `json:"border,omitempty"`
	MinWidth            string      `json:"min_width,omitempty"`
	Align               BlockAlign  `json:"align,omitempty"`
	Urgent              bool        `json:"urgent,omitempty"`
	Separator           bool        `json:"separator,omitempty"`
	SeparatorBlockWidth int         `json:"separator_block_width,omitempty"`
	Markup              BlockMarkup `json:"markup,omitempty"`
}

// A ClickEvent is sent by i3bar when the user clicks on a block
type ClickEvent struct {
	BlockIdentifier `yaml:",inline"`
	X               int `json:"x"`
	Y               int `json:"y"`
	Button          int `json:"button"`
}

// The StatusLine is just an list of Blocks
type StatusLine []Block

// StartIPC creates 2 chans to send statuslines and receive click events.
func StartIPC(h Header, w io.Writer, r io.Reader) (chan<- StatusLine, <-chan ClickEvent) {

	statuslines := make(chan StatusLine)

	go writeLoop(h, w, statuslines)
	statuslines <- []Block{}

	var events chan ClickEvent
	if h.ClickEvents {
		events = make(chan ClickEvent)
		go readLoop(r, events)
		<-events
	}

	return statuslines, events
}

func writeLoop(h Header, w io.Writer, s <-chan StatusLine) {
	writeJSON(w, &h)

	if _, err := w.Write([]byte("[\n")); err != nil {
		panic(err)
	}

	// The first statusline is sent for synchronisation
	<-s

	for l := range s {
		if err := writeJSON(w, &l); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing: %v\n", err)
		}
	}
}

func writeJSON(w io.Writer, data interface{}) (err error) {
	var b []byte
	if b, err = json.Marshal(&data); err == nil {
		if _, err = w.Write(b); err == nil {
			_, err = w.Write([]byte("\n"))
		}
	}
	return err
}

func readLoop(r io.Reader, e chan<- ClickEvent) {
	if _, err := fmt.Fscanf(r, "["); err != nil {
		panic(err)
	}

	s := bufio.NewScanner(r)

	// Synchronisation
	e <- ClickEvent{}

	for s.Scan() {
		var ce ClickEvent
		if err := json.Unmarshal(s.Bytes(), &ce); err == nil {
			e <- ce
		} else {
			fmt.Fprintf(os.Stderr, "Error reading %s: %v\n", s.Bytes(), err)
		}
	}

	if err := s.Err(); err != nil {
		panic(err)
	}
}

func (a BlockIdentifier) Matches(b BlockIdentifier) bool {
	return a.Name == b.Name && a.Name != "" && a.Instance == b.Instance && a.Instance != ""
}
