package modules

import (
	"bytes"
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
