package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mitchellh/cli"
)

type StreamCommand struct {
	client  http.Client
	output  io.Writer
	config  map[string]interface{}
	verbose bool
}

type plugin interface {
	init()
	run(*map[string]interface{})
}

var plugins []plugin

func register(f plugin) {
	plugins = append(plugins, f)
}

func (s *StreamCommand) Run(args []string) int {
	stream := flag.NewFlagSet("stream", flag.ContinueOnError)
	token := stream.String("token", "", "token for authenticating with api")
	jobID := stream.String("job-id", "", "id of job that was created")
	verbose := stream.Bool("verbose", false, "show request and response")
	if err := stream.Parse(args); err != nil {
		return -1
	}
	// read file, load to token
	if len(*token) == 0 {
		if tok, ok := s.config["token"].(string); ok && len(tok) > 0 {
			*token = tok
		} else {
			fmt.Println(s.Help())
			return -1
		}
	}
	s.verbose = *verbose
	req, err := http.NewRequest("GET", s.config["stream_url"].(string), nil)
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
	s.print("Response: %#v\n", resp)
	if resp.StatusCode == 401 {
		fmt.Println(`Invalid credentials`)
		return -1
	}
	if len(*jobID) == 0 {
		process(resp.Body)
	} else {
		readFromResponse(resp.Body, *jobID)
	}
	return 0
}

func process(r io.Reader) {
	dec := json.NewDecoder(r)
	m := make(map[string]interface{})
	for dec.More() {
		if err := dec.Decode(&m); err != nil {
			log.Fatal(err)
			continue
		}
		if len(m) == 0 {
			log.Println(m)
			continue
		}

		for _, p := range plugins {
			p.run(&m)
		}
	}
}

func (s *StreamCommand) print(pattern string, v interface{}) {
	if s.verbose {
		fmt.Fprintf(s.output, pattern, v)
	}
}

type msg struct {
	Origin origin `json:"origin"`
}

type origin struct {
	JobID string `json:"job_id"`
}

func readFromResponse(body io.ReadCloser, jobid string) {
	buf := bufio.NewReader(body)
	msg := &msg{}
	for {
		byts, _ := buf.ReadBytes('\n')

		if err := json.Unmarshal(byts, msg); err != nil {
			log.Println("error: ", err)
			break
		}
		if jobid == msg.Origin.JobID {
			fmt.Println(string(byts))
		}
	}
}

func (s *StreamCommand) Synopsis() string { return "receive stream from api" }

func (s *StreamCommand) Help() string {
	return `
Usage: 40fy-client stream -token=TOKEN [-job-id=JOBID]

 The "TOKEN" parameter is the token given to you by BinaryEdge, it is used as authentication.
 The "JOB-ID" parameter is optional, if it is present then the stream will be filtered for this specific job.
	`
}

func StreamCommandFactory() (cli.Command, error) {
	s := &StreamCommand{
		client: http.Client{
			Timeout: 5 * time.Second,
		},
		output: os.Stdout,
	}
	plugins = append(plugins, &printer{})
	if contents, err := GetConfigContents(config_file_name); err == nil {
		s.config = contents
	} else {
		s.config = DefaultConfig
	}
	for k, v := range DefaultConfig {
		if _, ok := s.config[k]; !ok {
			s.config[k] = v
		}
	}
	return s, nil
}
