package main

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/Adirelle/blocksgo/lib/framework"
	_ "github.com/Adirelle/blocksgo/lib/modules"
)

func main() {

	config := framework.Config{}

	if err := readConfig("blocksgo.yaml", &config); err != nil {
		panic(err)
	}

	app := framework.Application{Config: config}
	app.Run(os.Stdin, os.Stdout)
}

func readConfig(path string, config *framework.Config) (err error) {
	var (
		r   io.ReadCloser
		cfg []byte
	)

	if r, err = os.Open(path); err != nil {
		return
	}
	defer r.Close()

	if cfg, err = ioutil.ReadAll(r); err != nil {
		panic(err)
	}

	return framework.Configure(cfg, config)
}
