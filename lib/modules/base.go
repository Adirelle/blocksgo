package modules

import "time"

type PollingBlock struct {
	Interval time.Duration `yaml:"interval"`
	ticker   *time.Ticker
}

func (p PollingBlock) StartPolling() <-chan time.Time {
	p.ticker = time.NewTicker(p.Interval)
	return p.ticker.C
}

func (p PollingBlock) Stop() {
	p.ticker.Stop()
}
