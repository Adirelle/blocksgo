package modules

import (
	"fmt"
	"time"

	"github.com/Adirelle/blocksgo/lib/framework"
	"github.com/Adirelle/blocksgo/lib/ipc"
)

type Time struct {
	framework.BaseBlock `yaml:",inline"`
	PollingBlock        `yaml:",inline"`
	Format              string `yaml:"format"`
}

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

func init() {
	framework.RegistryBlockType("time", func() framework.Block { return &Time{} })
}
