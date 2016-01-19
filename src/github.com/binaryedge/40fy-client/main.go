package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
	"gopkg.in/BurntSushi/toml.v0"
)

const (
	CONFIG_PATH      = "CONFIG_PATH"
	config_home_path = ".binaryedge/"
	config_file_name = "config"
)

var (
	DefaultConfig = map[string]interface{}{
		"job_url":      `http://api.binaryedge.io/v1/tasks`,
		"stream_url":   `http://stream.api.binaryedge.io/v1/stream`,
		"firehose_url": `http://stream.api.binaryedge.io/v1/firehose`,
		"token":        "",
	}
)

func GetConfigContents(path string) (content map[string]interface{}, err error) {
	if _, err = toml.DecodeFile(path, &content); err != nil {
		return
	}
	return
}

func main() {
	c := cli.NewCLI("foreman", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"stream":     StreamCommandFactory,
		"firehose":   FirehoseCommandFactory,
		"create-job": CreateJobCommandFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
