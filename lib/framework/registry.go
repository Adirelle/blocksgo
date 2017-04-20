package framework

import "fmt"

type MakeFunc func() Block

var registry = make(map[string]MakeFunc)

// RegistryBlockType registers a block struct with a name
func RegistryBlockType(name string, make MakeFunc) {
	if _, exists := registry[name]; exists {
		panic(fmt.Sprintf("Block type %s already exists", name))
	}
	registry[name] = make
}

func newBlockInstance(name string) (Block, error) {
	make, exists := registry[name]
	if !exists {
		return nil, fmt.Errorf("Unknown block type %s", name)
	}

	return make(), nil
}
