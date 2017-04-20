package framework

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v2"
)

type GlobalConfig struct {
	Interval time.Duration `yaml:"interval"`
}

type preConfig struct {
	Global GlobalConfig
	Blocks []map[string]interface{} `yaml:"blocks"`
}

// Config contains the instances of blocks
type Config struct {
	Global GlobalConfig
	Blocks []Block
}

// Configure loads the configuration and instantiates thee blocks.
func Configure(yamlConfig []byte, conf *Config) (err error) {

	pre := preConfig{}
	if err = yaml.Unmarshal(yamlConfig, &pre); err != nil {
		return
	}

	conf.Global = pre.Global

	for _, bc := range pre.Blocks {
		namei, exists := bc["type"]
		if !exists {
			return fmt.Errorf("Block configuration must contain a 'type' key")
		}
		name, ok := namei.(string)
		if !ok {
			return fmt.Errorf("Block type must be a string")
		}

		var blk Block
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

		conf.Blocks = append(conf.Blocks, blk)
	}

	return
}
