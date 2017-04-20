package framework

import (
	"fmt"
	"reflect"
)

var registry = make(map[string]reflect.Type)

// RegistryBlockType registers a block struct with a name
func RegistryBlockType(name string, ptr *Block) {
	if _, exists := registry[name]; exists {
		panic(fmt.Sprintf("Block type %s already exists", name))
	}
	registry[name] = reflect.TypeOf(ptr).Elem()
}

func newBlockInstance(name string) (*Block, error) {
	t, exists := registry[name]
	if !exists {
		return nil, fmt.Errorf("Unknown block type %s", name)
	}

	b, ok := reflect.Zero(t).Interface().(Block)
	if !ok {
		return nil, fmt.Errorf("Could not cast %s to Block", b)
	}

	return &b, nil
}
