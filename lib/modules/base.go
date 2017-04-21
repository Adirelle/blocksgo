package modules

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

//============================================================

type PollFunc func(time.Time)

type PollingBlock struct {
	Interval time.Duration `yaml:"interval"`

	ticker *time.Ticker
}

func (p *PollingBlock) StartPolling(poll PollFunc) {
	p.ticker = time.NewTicker(p.Interval)
	go func() {
		poll(time.Now())
		for t := range p.ticker.C {
			poll(t)
		}
	}()
}

func (p *PollingBlock) Stop() {
	p.ticker.Stop()
}

//============================================================

type TemplatedBlock struct {
	Template string `yaml:"template"`

	tpl *template.Template
	buf bytes.Buffer
}

func (t *TemplatedBlock) ParseTemplate() {
	tplName := fmt.Sprintf("%p", t)
	t.tpl = template.Must(template.New(tplName).Parse(t.Template))
}

func (t *TemplatedBlock) Format(data interface{}) string {
	t.buf.Reset()
	if err := t.tpl.Execute(&t.buf, data); err != nil {
		return err.Error()
	}
	return t.buf.String()
}

//============================================================

type Unit string

const (
	BYTE Unit = "b"
	KILO Unit = "k"
	MEGA Unit = "M"
	GIGA Unit = "G"
)

var unitScales = map[Unit]float64{
	BYTE: 1.0,
	KILO: float64(1 << 10),
	MEGA: float64(1 << 20),
	GIGA: float64(1 << 30),
}

func (u Unit) Format(value interface{}) string {
	if s, ok := unitScales[u]; ok {
		if i, ok := value.(uint64); ok {
			return fmt.Sprintf("%.1f", float64(i)/s)
		} else {
			return fmt.Sprintf("Cannot cast %q into uint64", value)
		}
	} else {
		return fmt.Sprintf("Unit %q unknown", u)
	}
}
