package modules

import (
	"time"

	"golang.org/x/sys/unix"

	"github.com/Adirelle/blocksgo/lib/framework"
	"github.com/Adirelle/blocksgo/lib/ipc"
)

type Sysinfo struct {
	framework.BaseBlock `yaml:",inline"`
	PollingBlock        `yaml:",inline"`
	TemplatedBlock      `yaml:",inline"`
	Unit                Unit `yaml:"unit"`

	unix.Sysinfo_t
	Usedram  uint64
	Usedswap uint64
}

func init() {
	framework.RegistryBlockType("sysinfo", func() framework.Block {
		b := Sysinfo{}
		b.Name = "sysinfo"
		b.Interval = 1 * time.Second
		b.Unit = Giga
		return &b
	})
}

func (s *Sysinfo) Start() <-chan ipc.Block {
	s.ParseTemplate()
	s.StartPolling(func(_ time.Time) {
		b := s.MakeDefaultBlock()
		if err := unix.Sysinfo(&s.Sysinfo_t); err != nil {
			b.FullText = err.Error()
		} else {
			s.Usedram = s.Totalram - s.Freeram - s.Sharedram - s.Bufferram
			s.Usedswap = s.Totalswap - s.Freeswap
			b.FullText = s.Format(s)
		}
		s.Emit(b)
	})
	return s.BaseBlock.Start()
}

func (s *Sysinfo) Stop() {
	s.PollingBlock.Stop()
	s.BaseBlock.Stop()
}
