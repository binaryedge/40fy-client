package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/mitchellh/cli"
)

type FirehoseCommand struct {
	client  http.Client
	url     string
	output  io.Writer
	config  map[string]interface{}
	verbose bool
}

func (s *FirehoseCommand) Run(args []string) int {
	firehose := flag.NewFlagSet("firehose", flag.ContinueOnError)
	token := firehose.String("token", "", "")
	verbose := firehose.Bool("verbose", false, "show request and response")
	if err := firehose.Parse(args); err != nil {
		return -1
	}
	s.verbose = *verbose
	// read file, load to token
	if len(*token) == 0 {
		if file, err := ioutil.ReadFile("client_token"); err == nil {
			*token = string(file)
		}
	}
	if len(*token) == 0 {
		fmt.Println(s.Help())
		return -1
	}
	req, err := http.NewRequest("GET", s.url, nil)
	if err != nil {
		fmt.Println("Failed to connect ", err.Error())
		return -1
	}
	req.Header.Add("X-Token", *token)
	s.print("Request: %v\n", req)
	resp, err := s.client.Do(req)
	if err != nil {
		fmt.Println("Failed to connect ", err.Error())
		return -1
	}
	defer resp.Body.Close()
	s.print("Response: %v\n", resp)
	if resp.StatusCode == 401 {
		msg := `Invalid credentials`
		fmt.Println(msg)
		return -1
	}
	io.Copy(s.output, resp.Body)
	return 0
}

func (s *FirehoseCommand) print(pattern string, v interface{}) {
	if s.verbose {
		fmt.Fprintf(s.output, pattern, v)
	}
}

func (s *FirehoseCommand) Synopsis() string {
	return "Read JSON output from a stream with all content from the platform"
}

func (s *FirehoseCommand) Help() string {
	return `
Usage: be-cli firehose -token=TOKEN

 The "TOKEN" parameter is the token given to you by BinaryEdge, it is used as authentication.
	`
}

func FirehoseCommandFactory() (cli.Command, error) {
	f := &FirehoseCommand{
		client: http.Client{},
		output: os.Stdout,
	}
	if contents, err := GetConfigContents(config_file_name); err == nil {
		f.config = contents
	} else {
		f.config = DefaultConfig
	}
	return f, nil
}
