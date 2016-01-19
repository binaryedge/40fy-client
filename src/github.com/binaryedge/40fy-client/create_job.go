package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/mitchellh/cli"
)

type createJobCommand struct {
	config  map[string]interface{}
	output  io.Writer
	verbose bool
}

type jobRequest struct {
	Type        string    `json:"type"`
	Priority    bool      `json:"priority,omitempty"`
	Description string    `json:"description"`
	Options     []options `json:"options"`
}

type options struct {
	Worldscan bool      `json:"worldscan"`
	Ports     []portDef `json:"ports"`
	Targets   []string  `json:"targets,omitempty"`
}

type portDef struct {
	Port    int      `json:"port"`
	Sample  int      `json:"sample,omitempty"`
	Modules []string `json:"modules"`
}

func isCIDR(cidrs []string) bool {
	for _, cidr := range cidrs {
		if _, _, err := net.ParseCIDR(cidr); err != nil {
			return false
		}
	}
	return true
}

func isIP(ips []string) bool {
	for _, ip := range ips {
		if ip := net.ParseIP(ip); ip == nil {
			return false
		}
	}
	return true
}

func filterEmpty(s *[]string) {
	var ns []string
	for _, el := range *s {
		if len(el) > 0 {
			ns = append(ns, el)
		}
	}
	*s = ns
}

func (l *createJobCommand) Run(args []string) int {
	create := flag.NewFlagSet("create-job", flag.ContinueOnError)
	token := create.String("token", "", "authentication token")
	jobType := create.String("type", "scan", "type of scan")
	port := create.Int("port", 0, "Port to scan")
	sample := create.Int("sample", 0, "number of results needed")
	modules := create.String("modules", "", "modules of scan, example: ssh, ftp, service")
	targets := create.String("targets", "", "target of scan, example: 8.8.8.8")
	redirect := create.Bool("redirect", false, "flag shows stream of job created by command")
	verbose := create.Bool("verbose", false, "show request and response")
	if err := create.Parse(args); err != nil {
		return -1
	}
	l.verbose = *verbose

	if len(*token) == 0 {
		*token = l.config["token"].(string)
	}
	if len(*token) == 0 || len(*targets) == 0 {
		fmt.Println(l.Help())
		return -1
	}

	aTargets := strings.Split(*targets, ",")
	filterEmpty(&aTargets)
	if len(aTargets) == 0 {
		fmt.Println(l.Help())
		return -1
	}

	if !isIP(aTargets) && !isCIDR(aTargets) {
		fmt.Println(l.Help())
		return -1
	}

	aModules := strings.Split(*modules, ",")
	if len(aModules[0]) == 0 {
		aModules = []string{}
	}
	portdef := portDef{
		Port:    *port,
		Sample:  *sample,
		Modules: aModules,
	}
	opts := options{
		Worldscan: false,
		Ports:     []portDef{portdef},
		Targets:   aTargets,
	}
	job := jobRequest{
		Type: *jobType,
		//Priority: false,
		//Description: *desc,
		Options: []options{opts},
	}
	byts, err := json.Marshal(job)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to make request ", err)
		return -1
	}

	url := l.config["job_url"].(string)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(byts))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to make request ", err.Error())
		return -1
	}
	req.Header.Add("X-Token", *token)
	cl := http.Client{}
	resp, err := cl.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to make request ", err)
		return -1
	}
	bdy, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed reading response ", err)
		return -1
	}
	defer resp.Body.Close()
	s := struct {
		StreamURL string `json:"stream_url"`
		JobID     string `json:"job_id"`
		Message   string `json:"message"`
	}{}
	l.print("%v\n", string(bdy))
	if err = json.Unmarshal(bdy, &s); err != nil {
		fmt.Fprintf(os.Stderr, "Received invalid json %s", err.Error())
		return -1
	}
	if len(s.JobID) == 0 {
		fmt.Fprintf(os.Stderr, "Error in creating job, %s", s.Message)
		return -1
	}
	if *redirect {
		l.print("Redirecting to stream %s\n", l.config["stream_url"].(string))
		cmd, _ := StreamCommandFactory()
		return cmd.Run([]string{"-token=" + *token, "-job-id=" + s.JobID})
	} else {
		fmt.Println("You can connect to your stream with: ", s.StreamURL)
		fmt.Println("The identifier of the job is: ", s.JobID)
	}
	return 0
}

func (l *createJobCommand) print(pattern string, v interface{}) {
	if l.verbose {
		fmt.Fprintf(l.output, pattern, v)
	}
}

func (l *createJobCommand) Synopsis() string { return "Create a job in the platform" }

func (l *createJobCommand) Help() string {
	return `
Usage: 40fy-client create-job -token=TOKEN -targets=TARGETS -modules=MODULES -port=PORT [-redirect]

 The TOKEN parameter is the token given to you by BinaryEdge, it is used as authentication.
 The TARGETS parameter lists the hosts that will be targeted. Targets are a list of IPs or CIDRs.
 The MODULES parameter lists the modules used in the job.
 The PORT parameter is the port of the hosts that will be targeted in the job.
 The redirect is an optional flag that sets the command to retrieve the job output from the stream after creating the job.
	`
}

func CreateJobCommandFactory() (cli.Command, error) {
	j := &createJobCommand{
		output: os.Stdout,
	}
	if contents, err := GetConfigContents(config_file_name); err == nil {
		j.config = contents
	} else {
		j.config = DefaultConfig
	}
	return j, nil
}
