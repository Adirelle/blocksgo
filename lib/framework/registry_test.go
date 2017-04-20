package framework

import (
	"testing"

	"github.com/Adirelle/blocksgo/lib/ipc"
	"github.com/stretchr/testify/assert"
)

type Mock struct {
}

func (m *Mock) Start() <-chan ipc.Block {
	return nil
}

func (m *Mock) Stop() {
}

func TestUnknownBlockType(t *testing.T) {
	clearRegistry()
	_, err := newBlockInstance("bla")
	assert.Error(t, err)
}

func clearRegistry() {
	registry = make(map[string]MakeFunc)
}

func TestRegisterBlockType(t *testing.T) {
	clearRegistry()
	RegistryBlockType("mock", func() Block { return &Mock{} })
	b, err := newBlockInstance("mock")
	assert.Nil(t, err)
	assert.Implements(t, (*Block)(nil), b)
}

func TestRegisterDuplicateBlockType(t *testing.T) {
	clearRegistry()
	RegistryBlockType("mock", func() Block { return &Mock{} })
	assert.Panics(t, func() {
		RegistryBlockType("mock", func() Block { return &Mock{} })
	})
}
