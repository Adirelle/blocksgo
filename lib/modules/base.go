package modules

import "time"

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
