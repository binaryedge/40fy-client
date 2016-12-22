package main

import (
	"log"
)

type printer struct{}

func (p *printer) init() {}

func (p *printer) run(d *map[string]interface{}) {
	log.Println("PrintToStdout called with", d)
}
