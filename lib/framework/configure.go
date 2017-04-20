package framework

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type preConfig struct {
	Blocks []map[string]interface{} `yaml:"blocks"`
}

// Config contains the instances of blocks
type Config struct {
	Blocks []*Block
}

// Configure loads the configuration and instantiates thee blocks.
func Configure(yamlConfig []byte, c *Config) (err error) {

	pre := preConfig{}
	if err = yaml.Unmarshal(yamlConfig, &pre); err != nil {
		return
	}

	for _, bc := range pre.Blocks {
		namei, exists := bc["type"]
		if !exists {
			return fmt.Errorf("Block configuration must contain a 'type' key")
		}
		name, ok := namei.(string)
		if !ok {
			return fmt.Errorf("Block type must be a string")
		}

		var blk *Block
		if blk, err = newBlockInstance(name); err != nil {
			return
		}

		var yml []byte
		if yml, err = yaml.Marshal(bc); err != nil {
			return
		}

		if err = yaml.Unmarshal(yml, blk); err != nil {
			return
		}

		c.Blocks = append(c.Blocks, blk)
	}

	return
}
