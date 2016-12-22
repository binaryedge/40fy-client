package main

import (
	"log"
	"os"

	"gopkg.in/BurntSushi/toml.v0"
)

var (
	enc = toml.NewEncoder(os.Stdout)
)

func init() {
	register(&tomlPlugin{})
}

type tomlPlugin struct {
}

func (t *tomlPlugin) init() {}

func (t *tomlPlugin) run(d *map[string]interface{}) {
	//log.Println("toToml called with", d)

	if err := enc.Encode(&d); err != nil {
		log.Fatal(err)
	}
}
