# 40fy-client
Platform client

# Dependencies
* [Go](https://golang.org/dl/) , tested with Go 1.5.3

# Installation
* Clone this repo ```git clone git@github.com:binaryedge/40fy-client.git```
* Change folder to repo ```cd 40fy-client```
* Set this folder as GOPATH ```export GOPATH=$(pwd)``` (This step is necessary as part of Go [configuration](https://github.com/golang/go/wiki/GOPATH)
* Run install ```make all``` This step will create a binary in bin/

# Usage
* Token
  * A Token can be set with the flag ```--token=InsertYourToken``` when using the client, by creating a file with name ```config``` in the same directory as the 40fy-client binary with the content```token="InsertYourToken"```
* Mode Verbose
  * When using the ```--verbose``` debug messages are shown.
* Stream
  * ``` 40fy-client stream [--token=InsertYourToken] [--job-id=InsertYourJobID]```
  * When --job-id=ID is present, the stream will filter jobs with that ID otherwise will show everything from the user's stream. 
* Firehose
  * ``` 40fy-client firehose [--token=InsertYourToken] [--verbose]```
  * Shows jobs run by firehose.
* Create Job
  * ```40fy-client create-job [--token=InsertYourToken] -targets=Target -port=InsertPortToScan -sample=SampleSize -modules=ServiceToScan  [--verbose] [--redirect]```
    * The Targets are a comma separated set of ips, ```8.8.8.8,1.1.1.1```
    * The Sample size is the number of results necessary to satisfy a scan
    * The Modules are which modules to use in scan, example: http,service,ssl,ssh,vnc [link](https://github.com/binaryedge/api-publicdoc#supported-modules)
    * Port to scan.
