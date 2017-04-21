package modules

import (
	"time"

	"github.com/Adirelle/blocksgo/lib/framework"
	"github.com/Adirelle/blocksgo/lib/ipc"
)

type Time struct {
	framework.BaseBlock `yaml:",inline"`
	PollingBlock        `yaml:",inline"`
	Format              string `yaml:"format"`
}

func init() {
	framework.RegistryBlockType("time", func() framework.Block {
		b := Time{}
		b.Name = "time"
		b.Interval = 5 * time.Second
		b.Format = "2006-01-01 15:04:05 MST"
		return &b
	})
}

func (t *Time) Start() <-chan ipc.Block {
	t.StartPolling(func(now time.Time) {
		b := t.MakeDefaultBlock()
		b.FullText = now.Format(t.Format)
		t.Emit(b)
	})
	return t.BaseBlock.Start()
}

func (t *Time) Stop() {
	t.PollingBlock.Stop()
	t.BaseBlock.Stop()
}
